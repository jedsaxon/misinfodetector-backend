package handler

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MisinfoPayload struct {
	PostId         string `json:"postid"`
	Misinformation bool   `json:"misinfo"`
}

func (c *PostsController) HandleNewMisinfoReport(msg *amqp.Delivery) {
	var misinfoPayload MisinfoPayload
	if err := json.Unmarshal(msg.Body, &misinfoPayload); err != nil {
		log.Printf("recieved badly formatted payload from misinfo report: %v", err)
		return
	}

	post, err := c.dbs.FindPost(misinfoPayload.PostId)
	if err != nil {
		log.Printf("error finding post from misinfo report: %v", err)
		return
	}

	log.Printf("post with id %s report returned: %v", post.Id.String(), misinfoPayload.Misinformation)

	msg.Ack(false)
}
