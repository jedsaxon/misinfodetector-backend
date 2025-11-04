package models

import (
	"math/rand"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type (
	MisinfoState int64

	PostModel struct {
		Message       string       `json:"message"`
		Username      string       `json:"username"`
		SubmittedDate time.Time    `json:"date"`
		MisinfoState  MisinfoState `json:"misinfoState"`
	}

	PostModelId struct {
		Id            uuid.UUID    `json:"id"`
		Message       string       `json:"message"`
		Username      string       `json:"username"`
		SubmittedDate time.Time    `json:"date"`
		MisinfoState  MisinfoState `json:"misinfoState"`
	}
)

const (
	MisinfoStateFake MisinfoState = iota
	MisinfoStateTrue
	MisinfoStateNotChecked
)

// Creates a new post. Will strip spaces in username and message before creating it
func NewPost(message string, username string, misinfoState MisinfoState) *PostModel {
	submittedDate := time.Now()
	return &PostModel{
		Message:       strings.TrimSpace(message),
		Username:      strings.TrimSpace(username),
		MisinfoState:  misinfoState,
		SubmittedDate: submittedDate,
	}
}

func (p *PostModel) WithId(id uuid.UUID) *PostModelId {
	return &PostModelId{
		Id:            id,
		Message:       p.Message,
		Username:      p.Username,
		SubmittedDate: p.SubmittedDate,
		MisinfoState:  p.MisinfoState,
	}
}

func (p *PostModel) ValidatePost() map[string]string {
	errs := make(map[string]string)

	if len(p.Message) == 0 {
		errs["message"] = "Message cannot be empty"
	} else if len(p.Message) >= 256 {
		errs["message"] = "Message cannot contain more than 256 characters"
	}

	if len(p.Username) == 0 {
		errs["username"] = "Username cannot be empty"
	} else if len(p.Username) >= 64 {
		errs["username"] = "Username cannot be more than 64 characters"
	}

	return errs
}

func RandomPost() *PostModel {
	message := clampString(faker.Sentence(), 256)
	username := clampString(faker.Username(), 64)
	submittedDate := time.Now().AddDate(0, 0, -rand.Intn(60))
	containsMisinformation := rand.Intn(2)

	return &PostModel{
		Message:       message,
		Username:      username,
		SubmittedDate: submittedDate,
		MisinfoState:  MisinfoState(containsMisinformation),
	}
}

// clampString clamps the string to a specified length
func clampString(str string, max int) string {
	if len(str) > max {
		b := make([]byte, max)
		for i := range str {
			b[i] = str[i]
		}
		return string(b)
	}
	return str
}
