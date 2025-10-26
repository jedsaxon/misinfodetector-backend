package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type (
	PostModel struct {
		Message                string    `json:"message"`
		Username               string    `json:"username"`
		SubmittedDate          time.Time `json:"date"`
		ContainsMisinformation bool      `json:"potentialMisinformation"`
	}

	PostModelId struct {
		Id                     uuid.UUID `json:"id"`
		Message                string    `json:"message"`
		Username               string    `json:"username"`
		SubmittedDate          time.Time `json:"date"`
		ContainsMisinformation bool      `json:"potentialMisinformation"`
	}
)

// Creates a new post. Will strip spaces in username and message before creating it
func NewPost(message string, username string, containsMisinformation bool) *PostModel {
	submittedDate := time.Now()
	return &PostModel{
		Message:                strings.TrimSpace(message),
		Username:               strings.TrimSpace(username),
		ContainsMisinformation: containsMisinformation,
		SubmittedDate:          submittedDate,
	}
}


func (p *PostModel) WithId(id uuid.UUID) *PostModelId {
	return &PostModelId{
		Id:                     id,
		Message:                p.Message,
		Username:               p.Username,
		SubmittedDate:          p.SubmittedDate,
		ContainsMisinformation: p.ContainsMisinformation,
	}
}

func (p *PostModel) ValidatePost() map[string]string {
	errs := make(map[string]string)

	if len(p.Message) == 0 {
		errs["message"] = "Message cannot be empty"
	} else if len(p.Message) > 256 {
		errs["message"] = "Message cannot contain more than 256 characters"
	} 

	if len(p.Username) == 0 {
		errs["username"] = "Username cannot be empty"
	} else if len(p.Username) > 64 {
		errs["username"] = "Username cannot be more than 64 characters"
	}

	return errs
}

type Post struct {
	
}
