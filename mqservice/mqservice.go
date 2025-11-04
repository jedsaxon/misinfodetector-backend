package mqservice

import (
	"encoding/json"
	"log"
	"misinfodetector-backend/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	jsonAppType = "application/json"
	jsonEncType = "utf-8"
)

type (
	MqServiceConfig struct {
		Uri             string
		InputQueueName  string
		OutputQueueName string
		ConsumerName    string
	}

	MqService struct {
		conn        *amqp.Connection
		ch          *amqp.Channel
		inputQueue  *amqp.Queue
		outputQueue *amqp.Queue
		cfg         *MqServiceConfig
	}
)

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
		cfg:         cfg,
	}

	return mqService, closeMq, nil
}

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

// SubscribeToMisinfoOutput will consume outputQueue.Name, and execute callback with the
// delivery, in a goroutine. It is expected that callback will at some point call
// Delivery.Ack, to confirm with RabbitMQ that the action is complete.
func (mq *MqService) SubscribeToMisinfoOutput(callback func(*amqp.Delivery)) {
	msgs, err := mq.ch.Consume(mq.outputQueue.Name, mq.cfg.ConsumerName, false, false, false, false, nil)
	log.Printf("listening for new misinformation reports in rabbitmq/%s", mq.outputQueue.Name)
	if err != nil {
		log.Printf("error consuming output queue: %s", err)
		return
	}
	for msg := range msgs {
		go callback(&msg)
	}
}
