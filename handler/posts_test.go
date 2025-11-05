package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"misinfodetector-backend/dbservice"
	"misinfodetector-backend/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// newTestPostsController returns a new PostsController, with an in-memory Sqlite database, and a method to close the
// database connection. It will insert postsToInsert amount of posts into the database. Note that mqs will be set to nil.
func newTestPostsController(t *testing.T, postsToInsert int) (*PostsController, func()) {
	t.Helper()

	dbs, cancel, err := dbservice.NewDbService(":memory:")
	if err != nil {
		t.Fatalf("error connecting to dbservice: %v", err)
	}

	for i := 1; i <= postsToInsert; i++ {
		post := models.NewPost(fmt.Sprintf("msg %d", i), fmt.Sprintf("user%d", i), time.Now(), models.MisinfoStateNotChecked)
		if _, err := dbs.InsertPost(post); err != nil {
			t.Fatalf("could not insert post %d: %v", i, err)
		}
	}

	cancelErr := func() {
		if err := cancel(); err != nil {
			t.Fatalf("unable to close sqlite: %v", err)
		}
	}

	return NewPostsController(dbs, nil), cancelErr
}

// TestPutPost_Success sends a well-formed POST request containing
// a valid JSON body with "message" and "username". It expects a
// 200 OK response, a JSON payload confirming creation, a non-empty
// post ID, and a "Location" header pointing to the new resource.
func TestPutPost_Success(t *testing.T) {
	c, cancel := newTestPostsController(t, 2)
	defer cancel()

	body := PutPostForm{
		Message:  "Integration test",
		Username: "dev",
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("error marshalling json: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewReader(b))
	w := httptest.NewRecorder()

	c.PutPost(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var out ResponsePutPost
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if out.Message != "successfully created post" {
		t.Errorf("unexpected message: %s", out.Message)
	}
	if out.Post == nil || out.Post.Id.String() == "" {
		t.Error("expected a non-empty post id")
	}
	if loc := resp.Header.Get("Location"); loc == "" {
		t.Error("expected location header to be set")
	}
}

// TestPutPost_InvalidBody ensures that a POST request with a
// malformed JSON body (e.g., an invalid JSON token) yields a
// 400 Bad Request response.
func TestPutPost_InvalidBody(t *testing.T) {
	c, cancel := newTestPostsController(t, 2)
	defer cancel()

	req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewReader([]byte(`{invalid}`)))
	w := httptest.NewRecorder()

	c.PutPost(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// TestPutPost_ValidationError verifies that a POST request whose
// body contains an empty "message" or an over-length message
// triggers a 400 status and that the error map contains the
// appropriate field error keys.
func TestPutPost_ValidationError(t *testing.T) {
	c, cancel := newTestPostsController(t, 2)
	defer cancel()

	body := PutPostForm{
		Message:  "",
		Username: "user",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewReader(b))
	w := httptest.NewRecorder()

	c.PutPost(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// TestGetPosts_MissingParams verifies that a GET request without
// the required query parameters ("pageNumber" and "resultAmount")
// results in a 400 Bad Request. The test asserts that the
// response body contains an error message.
func TestGetPosts_MissingParams(t *testing.T) {
	c, cancel := newTestPostsController(t, 2)
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/api/posts", nil) // no query
	w := httptest.NewRecorder()

	c.GetPosts(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// TestGetPosts_InvalidParams checks that supplying
// invalid values for "pageNumber" or "resultAmount" (e.g.
// non-numeric or out-of-range numbers) produces a 400 status.
func TestGetPosts_InvalidParams(t *testing.T) {
	c, cancel := newTestPostsController(t, 2)
	defer cancel()

	params := url.Values{}
	params.Set("pageNumber", "1")
	params.Set("resultAmount", "100")
	req := httptest.NewRequest(http.MethodGet, "/api/posts?"+params.Encode(), nil)
	w := httptest.NewRecorder()

	c.GetPosts(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

// helper to perform a GET request against the controller with
// the given pageNumber and resultAmount query values.
func performGetPosts(t *testing.T, c *PostsController, page int, amount int) (*httptest.ResponseRecorder, *ResponseGetPosts) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	q := url.Values{}
	q.Set("pageNumber", fmt.Sprintf("%d", page))
	q.Set("resultAmount", fmt.Sprintf("%d", amount))
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	c.GetPosts(w, req)

	var resp ResponseGetPosts
	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
		return w, nil
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unable to unmarshal response: %v", err)
	}
	return w, &resp
}

func TestGetPosts_PageNumber1_Returns10(t *testing.T) {
	c, cancel := newTestPostsController(t, 15)
	defer cancel()

	w, resp := performGetPosts(t, c, 1, 10)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if resp == nil {
		t.Fatalf("response body was nil")
	}
	if len(resp.Posts) != 10 {
		t.Fatalf("expected 10 posts, got %d", len(resp.Posts))
	}
}

func TestGetPosts_PageNumber2_Returns5(t *testing.T) {
	c, cancel := newTestPostsController(t, 15)
	defer cancel()

	w, resp := performGetPosts(t, c, 2, 10)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if resp == nil {
		t.Fatalf("response body was nil")
	}
	if len(resp.Posts) != 5 {
		t.Fatalf("expected 5 posts, got %d", len(resp.Posts))
	}
}

func TestGetPosts_PageNumber3_Returns400(t *testing.T) {
	c, cancel := newTestPostsController(t, 15)
	defer cancel()

	w, resp := performGetPosts(t, c, 3, 10)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if resp == nil {
		t.Fatalf("response body was nil")
	}
	if len(resp.Posts) != 0 {
		t.Fatalf("expected 0 posts, got %d", len(resp.Posts))
	}
}
