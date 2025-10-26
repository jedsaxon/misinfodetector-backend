package dbservice

import (
	"time"

	"github.com/google/uuid"
)

type (
	PostModel struct {
		Message                string    `json:"message"`
		Username               string    `json:"username"`
		SubmittedDate          time.Time `json:"submittedDate"`
		ContainsMisinformation bool      `json:"containsMisinformation"`
	}

	PostModelId struct {
		Id                     uuid.UUID `json:"id"`
		Message                string    `json:"message"`
		Username               string    `json:"username"`
		SubmittedDate          time.Time `json:"submittedDate"`
		ContainsMisinformation bool      `json:"containsMisinformation"`
	}
)

func NewPost(message string, username string, containsMisinformation bool) *PostModel {
	submittedDate := time.Now()
	return &PostModel{
		Message:                message,
		Username:               username,
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
