package rabbitmq

import (
	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"
	godatabasessql "github.com/ralvarezdev/go-databases/sql"
	gojwtrabbitmqconsumer "github.com/ralvarezdev/go-jwt/rabbitmq/consumer"
	gojwtrabbitmqconsumersql "github.com/ralvarezdev/go-jwt/rabbitmq/consumer/sql"
	internalloader "github.com/ralvarezdev/rest-auth/internal/loader"
)

const (
	// EnvRabbitMQURL is the environment variable for the RabbitMQ URL
	EnvRabbitMQURL = "RABBITMQ_URL"

	// EnvRabbitMQTokenQueueName is the environment variable for the RabbitMQ token queue name
	EnvRabbitMQTokenQueueName = "RABBITMQ_TOKEN_QUEUE_NAME"

	// RabbitMQPollingPeriod is the RabbitMQ polling period for checking new messages
	RabbitMQPollingPeriod = 5 * time.Second

	// TokensMessagesConsumerChannelBufferSize is the buffer size for the tokens messages consumer channel
	TokensMessagesConsumerChannelBufferSize = 100
)

var (
	// RabbitMQURL is the RabbitMQ URL
	RabbitMQURL string

	// RabbitMQTokenQueueName is the RabbitMQ token queue name
	RabbitMQTokenQueueName string

	// RabbitMQConn is the RabbitMQ connection
	RabbitMQConn *amqp091.Connection

	// RabbitMQConsumer is the JWT RabbitMQ consumer
	RabbitMQConsumer gojwtrabbitmqconsumer.Consumer

	// RabbitMQService is the JWT RabbitMQ service
	RabbitMQService gojwtrabbitmqconsumersql.Service
)

// Load initializes the RabbitMQ constants
//
// Parameters:
//
//   - logger: The logger (optional, can be nil)
func Load(rabbitMQSQLHandler godatabasessql.Handler, logger *slog.Logger) {
	// Load the environment variables
	for env, dest := range map[string]*string{
		EnvRabbitMQURL:            &RabbitMQURL,
		EnvRabbitMQTokenQueueName: &RabbitMQTokenQueueName,
	} {
		if err := internalloader.Loader.LoadVariable(
			env,
			dest,
		); err != nil {
			panic(err)
		}
	}

	// Create a RabbitMQ connection
	conn, err := amqp091.Dial(RabbitMQURL)
	if err != nil {
		panic(err)
	}
	RabbitMQConn = conn

	// Create the JWT RabbitMQ consumer
	rabbitMQConsumer, err := gojwtrabbitmqconsumer.NewDefaultConsumer(
		RabbitMQConn,
		RabbitMQTokenQueueName,
		RabbitMQPollingPeriod,
		TokensMessagesConsumerChannelBufferSize,
		logger,
	)
	if err != nil {
		panic(err)
	}
	RabbitMQConsumer = rabbitMQConsumer

	// Create the JWT RabbitMQ service
	rabbitMQService, err := gojwtrabbitmqconsumersql.NewDefaultService(
		rabbitMQSQLHandler,
		RabbitMQConsumer,
		logger,
	)
	if err != nil {
		panic(err)
	}
	RabbitMQService = rabbitMQService
}
