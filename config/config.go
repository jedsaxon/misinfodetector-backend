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
}

const (
	defaultSqliteDsn           = ":memory:"
	defaultListenAddress       = "127.0.0.1:5000"
	defaultRabbitMqUri         = "amqp://guest:guest@localhost:5672/"
	defaultPostOutputQueueName = "misinfo/output"
	defaultPostInputQueueName  = "misinfo/input"
)

func NewDefaultConfig() *Config {
	return &Config{
		SqliteDsn:               defaultSqliteDsn,
		RabbitMqUri:             defaultRabbitMqUri,
		RabbitMqInputQueueName:  defaultPostInputQueueName,
		RabbitMqOutputQueueName: defaultPostOutputQueueName,
		ListenAddres:            defaultListenAddress,
	}
}

func (c *Config) PopulateFromArgs() {
	var sqliteDsn = flag.String("sqlite", "", "where the sqlite database should be stored (default 127.0.0.1:5000)")
	var rabbitMqUri = flag.String("rabbitmq", "", "rabbitmq connection string (default amqp://guest:guest@localhost:5672/)")
	var inputQueueName = flag.String("inputqueue", "", "rabbitmq input queue name (default misinfo/input)")
	var outputQueueName = flag.String("outputqueue", "", "rabbitmq input queue name (default misinfo/output)")
	var listenAddress = flag.String("listen", "", "where this program should listen for api requests (default :memory:)")

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
}

func (c *Config) PopulateFromEnv() {
	myEnv, err := godotenv.Read()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf(".env not found, won't load environment variables")
		} else {
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

	if inputQueuName, ok := myEnv["INPUT_QUEUE_NAME"]; ok {
		c.RabbitMqInputQueueName = inputQueuName
	}

	if outputQueuName, ok := myEnv["OUTPUT_QUEUE_NAME"]; ok {
		c.RabbitMqInputQueueName = outputQueuName
	}
}
