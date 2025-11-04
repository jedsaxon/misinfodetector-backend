package mqservice

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type MqService struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue *amqp.Queue
}

func NewMqService(rabbitMQUri string) (*MqService, func() error, error) {
	connection, err := amqp.Dial(rabbitMQUri)
	if err != nil {
		return nil, nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, nil, err
	}

	queue, err := channel.QueueDeclare(
		"misinfo", 
		false,   
		false,   
		false,   
		false,   
		nil,     
	)

	closeMq := func() error {
		if err := connection.Close(); err != nil {
			return err
		}
		if err := channel.Close(); err != nil {
			return err
		}
		return nil
	}

	mqService := &MqService{
		conn: connection,
		ch:   channel,
		queue: &queue,
	}

	return mqService, closeMq, nil
}

