package handler

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *PostsController) HandleNewMisinfoReport(msg *amqp.Delivery) {
	log.Printf("recieved report: ", msg.Body)
	msg.Ack(false)
}
