package handler

import (
	"fmt"
	"log"
	"math"
	"misinfodetector-backend/dbservice"
	"misinfodetector-backend/handler/util"
	"misinfodetector-backend/models"
	"misinfodetector-backend/mqservice"
	"misinfodetector-backend/validation"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

type (
	PostsController struct {
		dbs *dbservice.DbService
		mqs *mqservice.MqService
	}

	PutPostForm struct {
		Message  string `json:"message"`
		Username string `json:"username"`
	}

	PutRandomPostsForm struct {
		Amount int `json:"amount"`
	}

	ResponsePutPost struct {
		Message string              `json:"message"`
		Post    *models.PostModelId `json:"post"`
	}

	ResponseFindPost struct {
		Message string              `json:"message"`
		Post    *models.PostModelId `json:"post"`
	}

	ResponseRandomPosts struct {
		Message string `json:"message"`
		Amount  int    `json:"amount"`
	}

	ResponseGetPosts struct {
		Message   string                `json:"message"`
		Posts     []*models.PostModelId `json:"posts"`
		PageCount int64                 `json:"pages"`
	}
)

func NewPostsController(dbs *dbservice.DbService, mqs *mqservice.MqService) *PostsController {
	return &PostsController{
		dbs: dbs,
		mqs: mqs,
	}
}

func (c *PostsController) GetPosts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageNumber, resultAmount, errs := validation.ValidateGetPostsRequest(query)
	if len(errs) > 0 {
		New400Response(errs).RespondToFatal(w)
		return
	}

	postCount, err := c.dbs.GetPostCount()
	if err != nil {
		log.Printf("unable to get post count: %v", err)
		New500Response().RespondToFatal(w)
		return
	}

	posts, err := c.dbs.GetPosts(pageNumber, resultAmount)
	if err != nil {
		log.Printf("unable to get posts: %v", err)
		New500Response().RespondToFatal(w)
		return
	}

	pageCount := math.Ceil(float64(postCount) / float64(resultAmount))
	WriteJsonFatal(http.StatusOK, w, &ResponseGetPosts{
		Message:   fmt.Sprintf("%d posts found", len(posts)),
		Posts:     posts,
		PageCount: int64(pageCount),
	})
}

func (c *PostsController) PutPosts(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("posts")
	if err != nil {
		errs := make(map[string]string, 0)
		errs["posts"] = "file was not found"
		New400Response(errs).RespondToFatal(w)
	}
	defer f.Close()

	err = c.dbs.ImportPosts(f)
	if err != nil {
		log.Printf("unable to insert posts: %v", err)
		New500Response().RespondToFatal(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *PostsController) GetSpecificPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Printf("expected 'id' from mux.Vars(), got nothing")
		New500Response().RespondToFatal(w)
		return
	}

	post, err := c.dbs.FindPost(id)
	if err != nil {
		log.Printf("unable to find post: %v", err)
		New500Response().RespondToFatal(w)
		return
	}

	if post == nil {
		errs := make(map[string]string)
		errs["id"] = "Post with the given ID could not be found"
		New400Response(errs).RespondToFatal(w)
		return
	}

	WriteJsonFatal(http.StatusOK, w, &ResponseFindPost{
		Message: "found post with the given ID",
		Post:    post,
	})
}

func (c *PostsController) UploadPost(w http.ResponseWriter, r *http.Request) {
	var body PutPostForm
	if err := util.UnmarshalJsonReader(r.Body, &body); err != nil {
		NewCustomResponse(http.StatusBadRequest, "malformed body").RespondTo(w)
		log.Printf("unable to unmarshal body: %v", err)
		return
	}
	log.Printf("hi")

	submittedDate := time.Now().UTC()
	post := models.NewPost(body.Message, body.Username, submittedDate)
	if errs := post.ValidatePost(); len(errs) > 0 {
		New400Response(errs).RespondToFatal(w)
		return
	}

	postWithId, err := c.dbs.InsertPost(post)
	if err != nil {
		New500Response().RespondToFatal(w)
		log.Printf("error inserting post: %v", err)
		return
	}

	if c.mqs != nil {
		log.Printf("publishing new post to rabbitmq connection")
		err = c.mqs.PublishNewPost(postWithId)
		if err != nil {
			New500Response().RespondToFatal(w)
			log.Printf("error inserting post: %v", err)
			return
		}
	} else {
		log.Printf("no rabbitmq connection found - will not publish new post")
	}

	createdUrl, err := url.JoinPath(r.Host, "api", "posts", postWithId.Id.String())
	if err != nil {
		New500Response().RespondToFatal(w)
		log.Fatalf("unable to generate created URL: %v", err)
		return
	}
	w.Header().Add("location", createdUrl)

	WriteJsonFatal(http.StatusOK, w, &ResponsePutPost{
		Message: "successfully created post",
		Post:    postWithId,
	})
}

func (c *PostsController) PutRandomPosts(w http.ResponseWriter, r *http.Request) {
	var body PutRandomPostsForm
	if err := util.UnmarshalJsonReader(r.Body, &body); err != nil {
		NewCustomResponse(http.StatusBadRequest, "malformed body").RespondTo(w)
		log.Printf("unable to unmarshal body: %v", err)
		return
	}

	if err := validation.ValidateRandomAmount(body.Amount); err != nil {
		errs := make(map[string]string)
		errs["amount"] = err.Error()
		New400Response(errs).RespondToFatal(w)
		return
	}

	for _ = range body.Amount {
		p := models.RandomPost()
		_, err := c.dbs.InsertPost(p)
		if err != nil {
			log.Printf("unable to create post: %v", err)
			New500Response().RespondToFatal(w)
			return
		}
	}

	WriteJsonFatal(http.StatusOK, w, &ResponseRandomPosts{
		Message: fmt.Sprintf("successfully created %d random posts", body.Amount),
		Amount:  body.Amount,
	})
}

func DeletePosts(w http.ResponseWriter, r *http.Request, db *dbservice.DbService) {

}
