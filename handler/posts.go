package handler

import (
	"fmt"
	"log"
	"math"
	"misinfodetector-backend/dbservice"
	"misinfodetector-backend/handler/util"
	"misinfodetector-backend/models"
	"misinfodetector-backend/validation"
	"net/http"
	"net/url"
)

type (
	PostsController struct {
		dbs *dbservice.DbService
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

	ResponseRandomPosts struct {
		Message string `json:"message"`
		Amount  int    `json:"amount"`
	}

	ResponseGetPosts struct {
		Message   string               `json:"message"`
		Posts     []models.PostModelId `json:"posts"`
		PageCount int64                `json:"pages"`
	}
)

func NewPostsController(dbs *dbservice.DbService) *PostsController {
	return &PostsController{
		dbs: dbs,
	}
}

func (c *PostsController) GetPosts(w http.ResponseWriter, r *http.Request) {
	const (
		pageNumberQueryName   = "pageNumber"
		resultAmountQueryName = "resultAmount"
	)

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

func (c *PostsController) PutPost(w http.ResponseWriter, r *http.Request) {
	var body PutPostForm
	if err := util.UnmarshalJsonReader(r.Body, &body); err != nil {
		NewCustomResponse(http.StatusBadRequest, "malformed body").RespondTo(w)
		log.Printf("unable to unmarshal body: %v", err)
		return
	}

	post := models.NewPost(body.Message, body.Username, false)
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
		Amount: body.Amount,
	})
}

func DeletePosts(w http.ResponseWriter, r *http.Request, db *dbservice.DbService) {

}
