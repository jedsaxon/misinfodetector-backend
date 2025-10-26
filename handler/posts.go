package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"misinfodetector-backend/dbservice"
	"misinfodetector-backend/handler/validation"
	"misinfodetector-backend/util"
	"net/http"
	"net/url"
)

type (
	PutPostForm struct {
		Message  string `json:"message"`
		Username string `json:"username"`
	}

	ResponseGetPosts struct {
		Message   string                  `json:"message"`
		Posts     []dbservice.PostModelId `json:"posts"`
		PageCount int                     `json:"pages"`
	}
)

func GetPosts(w http.ResponseWriter, r *http.Request, db *dbservice.DbService) {
	const (
		pageNumberQueryName   = "pageNumber"
		resultAmountQueryName = "resultAmount"
	)
	w.Header().Add("Content-Type", "application/json")

	query := r.URL.Query()
	pageNumber, resultAmount, errs := validation.ValidateGetPostsRequest(query)
	if len(errs) > 0 {
		util.New400Response(errs).RespondToFatal(w)
		return
	}

	posts, err := db.GetPosts(pageNumber+1, resultAmount)
	if err != nil {
		log.Printf("unable to get posts: %v", err)
		return
	}

	responseJson, err := json.Marshal(&ResponseGetPosts{
		Message:   fmt.Sprintf("%d posts found", len(posts)),
		Posts:     posts,
		PageCount: len(posts),
	})
	if err != nil {
		util.New500Response().RespondToFatal(w)
		log.Printf("unable to marshal response: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(responseJson))
	if err != nil {
		log.Printf("unable to write response to user: %v", err)
		return
	}
}

func PutPost(w http.ResponseWriter, r *http.Request, db *dbservice.DbService) {
	bodyBytes, err := io.ReadAll(r.Body)
	w.Header().Add("Content-Type", "application/json")

	if err != nil {
		util.New500Response().RespondTo(w)
		log.Printf("error reading body: %v", err)
		return
	}

	var body PutPostForm
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		util.NewCustomResponse(http.StatusBadRequest, "malformed body").RespondTo(w)
		log.Printf("unable to unmarshal body: %v", err)
		return
	}

	post := dbservice.NewPost(body.Message, body.Username, false)
	postWithId, err := db.InsertPost(post)
	if err != nil {
		util.New500Response().RespondTo(w)
		log.Printf("error inserting post: %v", err)
		return
	}

	response := struct {
		Message string                 `json:"message"`
		Post    *dbservice.PostModelId `json:"post"`
	}{
		Message: "successfully created post",
		Post:    postWithId,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		util.New500Response().RespondTo(w)
		log.Printf("error marshalling response: %v", err)
		return
	}
	createdUrl, err := url.JoinPath(r.Host, "api", "posts", postWithId.Id.String())
	if err != nil {
		util.New500Response().RespondTo(w)
		log.Fatalf("unable to generate created URL: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("location", createdUrl)

	if _, err = w.Write(responseJson); err != nil {
		log.Printf("error writing to socket: %v", err)
	}
}

func DeletePosts(w http.ResponseWriter, r *http.Request, db *dbservice.DbService) {

}
