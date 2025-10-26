package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"misinfodetector-backend/dbservice"
	"misinfodetector-backend/util"
	"net/http"
	"strconv"
	"strings"
)

func GetPosts(w http.ResponseWriter, r *http.Request, db *dbservice.DbService) {
	const (
		pageNumberQueryName   = "pageNumber"
		resultAmountQueryName = "resultAmount"
	)
	query := r.URL.Query()

	pageNumberQuery := strings.TrimSpace(query.Get(pageNumberQueryName))
	resultAmountQuery := strings.TrimSpace(query.Get(resultAmountQueryName))

	errors := make(map[string]string)
	if pageNumberQuery == "" {
		errors[pageNumberQueryName] = "you must provide a pageNumber"
	}
	pageNumber, err := strconv.ParseInt(pageNumberQuery, 10, 64)
	if err != nil {
		errors[pageNumberQueryName] = err.Error()
	} else if pageNumber <= 0 {
		errors[pageNumberQueryName] = "page number cannot be bellow 0"
	}

	if resultAmountQuery == "" {
		errors[resultAmountQueryName] = "you must provide a resultAmount"
	}
	resultAmount, err := strconv.ParseInt(resultAmountQuery, 10, 64)
	if err != nil {
		errors[resultAmountQueryName] = err.Error()
	} else if resultAmount <= 0 {
		errors[resultAmountQueryName] = "resultAmount must be greater than 0"
	} else if resultAmount >= 50 {
		errors[resultAmountQueryName] = "resultAmount must be less than 50"
	}

	if len(errors) > 0 {
		util.New400Response(errors).RespondToFatal(w)
		return
	}

	posts, err := db.GetPosts(pageNumber, resultAmount)
	if err != nil {
		log.Fatalf("unable to get posts: %v", err)
	}

	response := struct {
		Message string                `json:"message"`
		Posts   []dbservice.PostModel `json:"posts"`
	}{
		Message: fmt.Sprintf("%d posts found", len(posts)),
		Posts:   posts,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		util.New500Response().RespondToFatal(w)
		log.Fatalf("unable to marshal response: %v", err)
	}

	_, err = w.Write([]byte(responseJson))
	if err != nil {
		log.Fatalf("unable to write response to user: %v", err)
	}
}

func PutPosts(w http.ResponseWriter, r *http.Request, db *dbservice.DbService) {

}

func DeletePosts(w http.ResponseWriter, r *http.Request, db *dbservice.DbService) {

}
