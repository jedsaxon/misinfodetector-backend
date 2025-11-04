package mqservice

import (
	"encoding/json"
	"misinfodetector-backend/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MqServiceConfig struct {
	Uri             string
	InputQueueName  string
	OutputQueueName string
}

type MqService struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	inputQueue  *amqp.Queue
	outputQueue *amqp.Queue
}

func NewMqService(cfg *MqServiceConfig) (*MqService, func() error, error) {
	connection, err := amqp.Dial(cfg.Uri)
	if err != nil {
		return nil, nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, nil, err
	}

	inputQueue, err := channel.QueueDeclare(
		cfg.InputQueueName,
		false,
		false,
		false,
		false,
		nil,
	)

	outputQueue, err := channel.QueueDeclare(
		cfg.OutputQueueName,
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
		conn:        connection,
		ch:          channel,
		inputQueue:  &inputQueue,
		outputQueue: &outputQueue,
	}

	return mqService, closeMq, nil
}

const (
	jsonAppType = "application/json"
	jsonEncType = "utf-8"
)

// PublishNewPost publishes a post to the message queue
func (mq *MqService) PublishNewPost(p *models.PostModelId) error {
	body, err := json.Marshal(p)
	if err != nil {
		return err
	}

	payload := amqp.Publishing{
		ContentType:     jsonAppType,
		ContentEncoding: jsonEncType,
		Body:            body,
	}

	err = mq.ch.Publish("", mq.inputQueue.Name, false, false, payload)
	if err != nil {
		return err
	}

	return nil
}
