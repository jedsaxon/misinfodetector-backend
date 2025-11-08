package handler

import (
	"encoding/json"
	"log"
	"misinfodetector-backend/models"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type optionalString string

type MisinfoPayload struct {
	PostId           uuid.UUID           `json:"post_id"`
	Misinformation   models.MisinfoState `json:"misinfo_state"`
	DateSubmittedUTC string              `json:"date_submitted"`
}

func (c *PostsController) HandleNewMisinfoReport(msg *amqp.Delivery) {
	var misinfoPayload MisinfoPayload

	if err := json.Unmarshal(msg.Body, &misinfoPayload); err != nil {
		msg.Ack(false)
		log.Printf("recieved badly formatted payload from misinfo report: %v", err)
		return
	}

	if misinfoPayload.PostId == uuid.Nil {
		msg.Ack(false)
		log.Printf("recieved badly formatted payloadf rom misinfo report: null/nonexistent id")
		return
	}

	post, err := c.dbs.FindPost(misinfoPayload.PostId.String())
	if err != nil {
		msg.Ack(false)
		log.Printf("error finding post from misinfo report: %v", err)
		return
	}
	if post == nil {
		msg.Ack(false)
		log.Printf("error finding post from misinfo report: post with id [%s] not found", misinfoPayload.PostId.String())
		return
	}

	updatedPost := models.NewPost(post.Message, post.Username, post.SubmittedDateUTC)
	updatedPost.AttachReportToPost(post.MisinfoReport.State, post.MisinfoReport.Confidence, post.SubmittedDateUTC)
	upd, err := c.dbs.UpdatePost(post, updatedPost)
	if err != nil {
		log.Printf("error updating post from misinfo payload: %v", err)
		return
	} else if upd <= 0 {
		log.Printf("failed to update post: no records updated in database")
		return
	}

	log.Printf("post with id %s report returned: %v", post.Id.String(), misinfoPayload.Misinformation)

	msg.Ack(false)
}
