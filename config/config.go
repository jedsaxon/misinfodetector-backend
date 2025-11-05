package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SqliteDsn               string
	ListenAddres            string
	RabbitMqUri             string
	RabbitMqOutputQueueName string
	RabbitMqInputQueueName  string
	RabbitMqConsumerName    string
}

const (
	defaultSqliteDsn            = ":memory:"
	defaultListenAddress        = "127.0.0.1:5000"
	defaultRabbitMqUri          = "amqp://guest:guest@localhost:5672/"
	defaultRabbitMqConsumerName = "api"
	defaultPostOutputQueueName  = "misinfo/output"
	defaultPostInputQueueName   = "misinfo/input"
)

func NewDefaultConfig() *Config {
	return &Config{
		SqliteDsn:               defaultSqliteDsn,
		RabbitMqUri:             defaultRabbitMqUri,
		RabbitMqInputQueueName:  defaultPostInputQueueName,
		RabbitMqOutputQueueName: defaultPostOutputQueueName,
		RabbitMqConsumerName:    defaultRabbitMqConsumerName,
		ListenAddres:            defaultListenAddress,
	}
}

func (c *Config) PopulateFromArgs() {
	var sqliteDsn = flag.String("sqlite", "", "where the sqlite database should be stored (default 127.0.0.1:5000)")
	var rabbitMqUri = flag.String("rabbitmq", "", "rabbitmq connection string (default amqp://guest:guest@localhost:5672/)")
	var inputQueueName = flag.String("inputqueue", "", "rabbitmq input queue name (default misinfo/input)")
	var outputQueueName = flag.String("outputqueue", "", "rabbitmq input queue name (default misinfo/output)")
	var listenAddress = flag.String("listen", "", "where this program should listen for api requests (default :memory:)")
	var consumerName = flag.String("rabbitmq-name", "", "name to give the rabbitmq consumer")

	flag.Parse()

	if *sqliteDsn != "" {
		c.SqliteDsn = *sqliteDsn
	}

	if *listenAddress != "" {
		c.ListenAddres = *listenAddress
	}

	if *rabbitMqUri != "" {
		c.RabbitMqUri = *rabbitMqUri
	}

	if *inputQueueName != "" {
		c.RabbitMqInputQueueName = *inputQueueName
	}

	if *outputQueueName != "" {
		c.RabbitMqOutputQueueName = *outputQueueName
	}

	if *consumerName != "" {
		c.RabbitMqConsumerName = *consumerName
	}
}

func (c *Config) PopulateFromEnv() {
	myEnv, err := godotenv.Read()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalf("unable to read environment variables: %v", err)
		}
	}

	if sqliteDsn, ok := myEnv["SQLITE_DSN"]; ok {
		c.SqliteDsn = sqliteDsn
	}

	if listenAddress, ok := myEnv["LISTEN_ADDRESS"]; ok {
		c.ListenAddres = listenAddress
	}

	if rabbitMqUri, ok := myEnv["RABBITMQ_URI"]; ok {
		c.RabbitMqUri = rabbitMqUri
	}

	if consumerName, ok := myEnv["RABBITMQ_CONSUMER_NAME"]; ok {
		c.RabbitMqConsumerName = consumerName
	}

	if inputQueuName, ok := myEnv["RABBITMQ_INPUT_QUEUE_NAME"]; ok {
		c.RabbitMqInputQueueName = inputQueuName
	}

	if outputQueuName, ok := myEnv["RABBITMQ_OUTPUT_QUEUE_NAME"]; ok {
		c.RabbitMqInputQueueName = outputQueuName
	}
}
