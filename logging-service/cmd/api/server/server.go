package server

import (
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/handlers"
	"github.com/Salladin95/rmqtools"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

type App interface {
	Start()
}

type app struct {
	rabbit *amqp091.Connection
	mongo  *mongo.Client
	cfg    *config.Config
}

func NewApp(
	rabbit *amqp091.Connection,
	mongo *mongo.Client,
	cfg *config.Config,
) App {
	return &app{rabbit: rabbit, mongo: mongo, cfg: cfg}
}

func (a *app) Start() {
	topics := []string{constants.LogCommand}
	handlers := handlers.NewLoggingHandlers(a.mongo)
	// Create a new consumer for the specified AMQP exchange and queue.
	consumer, err := rmqtools.NewConsumer(
		a.rabbit,
		constants.AmqpExchange,
		constants.AmqpQueue,
	)
	if err != nil {
		log.Printf("failed to create consumer - %v\n", err) // Log error if connection fails
		os.Exit(1)                                          // Exit program if connection fails
		return
	}

	// Start listening for messages on the specified topics and invoke the message handler.
	consumer.Listen(topics, handlers.Log)
}
