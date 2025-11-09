package handler

import (
	"log"
	"net/http"
)

func (c *PostsController) GetTnseEmbeddings(w http.ResponseWriter, r *http.Request) {
	records, err := c.dbs.GetAllTnseEmbeddings()
	if err != nil {
		log.Printf("error getting all tnse embeddings: %v", err)
		New500Response().RespondToFatal(w)
		return
	}

	WriteJsonFatal(http.StatusOK, w, records)
}

func (c *PostsController) PutTnseEmbeddings(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("embeddings")
	if err != nil {
		errs := make(map[string]string, 0)
		errs["embeddings"] = "file was not found"
		New400Response(errs).RespondToFatal(w)
		return
	}
	defer f.Close()

	err = c.dbs.ImportTnseEmbeddings(f)
	if err != nil {
		log.Printf("unable to import tnse embeddings: %v", err)
		New500Response().RespondToFatal(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *PostsController) GetTopicActivities(w http.ResponseWriter, r *http.Request) {
	records, err := c.dbs.GetAllTopicActivityRecords()
	if err != nil {
		log.Printf("error getting all topic activities: %v", err)
		New500Response().RespondToFatal(w)
		return
	}

	WriteJsonFatal(http.StatusOK, w, records)
}

func (c *PostsController) PutTopicActivities(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("topics")
	if err != nil {
		errs := make(map[string]string, 0)
		errs["topics"] = "file was not found"
		New400Response(errs).RespondToFatal(w)
		return
	}
	defer f.Close()

	err = c.dbs.ImportTopicActivities(f)
	if err != nil {
		log.Printf("unable to import topic activities: %v", err)
		New500Response().RespondToFatal(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
