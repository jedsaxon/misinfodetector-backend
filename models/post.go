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
		Message          string    `json:"message"`
		Username         string    `json:"username"`
		SubmittedDateUTC time.Time `json:"date"`
		// MisinfoReport contains details about whether this post is misinformation.
		// Will be nil if a report has not been made
		MisinfoReport *MisinformationReport `json:"misinfo_report,omitempty"`
	}

	MisinformationReport struct {
		State            MisinfoState `json:"state"`
		Confidence       float32      `json:"confidence"`
		SubmittedDateUtc time.Time    `json:"submitted_date"`
	}

	PostModelId struct {
		Id               uuid.UUID `json:"id"`
		Message          string    `json:"message"`
		Username         string    `json:"username"`
		SubmittedDateUTC time.Time `json:"date"`
		// MisinfoReport contains details about whether this post is misinformation.
		// Will be nil if a report has not been made
		MisinfoReport *MisinformationReport `json:"misinfo_report"`
	}
)

const (
	MisinfoStateFake MisinfoState = iota
	MisinfoStateTrue
)

// Creates a new post. Will strip spaces in username and message before creating it
func NewPost(message string, username string, submittedDateUtc time.Time) *PostModel {
	return &PostModel{
		Message:          strings.TrimSpace(message),
		Username:         strings.TrimSpace(username),
		MisinfoReport:    nil,
		SubmittedDateUTC: submittedDateUtc,
	}
}

func (p *PostModel) AttachReportToPost(state MisinfoState, confidence float32, submittedDateUtc time.Time) {
	p.MisinfoReport = &MisinformationReport{
		State:            state,
		Confidence:       confidence,
		SubmittedDateUtc: submittedDateUtc,
	}
}

func (p *PostModelId) AttachReportToPost(state MisinfoState, confidence float32, submittedDateUtc time.Time) {
	p.MisinfoReport = &MisinformationReport{
		State:            state,
		Confidence:       confidence,
		SubmittedDateUtc: submittedDateUtc,
	}
}

func (p *PostModel) WithId(id uuid.UUID) *PostModelId {
	var duplicateReport *MisinformationReport = nil
	if p.MisinfoReport != nil {
		duplicateReport = &MisinformationReport{
			State:            p.MisinfoReport.State,
			Confidence:       p.MisinfoReport.Confidence,
			SubmittedDateUtc: p.MisinfoReport.SubmittedDateUtc,
		}
	}

	return &PostModelId{
		Id:               id,
		Message:          p.Message,
		Username:         p.Username,
		SubmittedDateUTC: p.SubmittedDateUTC,
		MisinfoReport:    duplicateReport,
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
	submittedDate := time.Now().UTC().AddDate(0, 0, -rand.Intn(60))

	return &PostModel{
		Message:          message,
		Username:         username,
		SubmittedDateUTC: submittedDate,
		MisinfoReport:    nil,
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
