package main

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/mail-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/mail-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/mail-service/cmd/api/handlers"
	"github.com/Salladin95/card-quizzler-microservices/mail-service/cmd/api/mail"
	"github.com/Salladin95/rmqtools"
	"log"
	"os"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Connect to RabbitMQ server using the provided URL.
	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.RabbitUrl)
	if err != nil {
		log.Println(err) // Log error if RabbitMQ connection fails
		os.Exit(1)       // Exit program if RabbitMQ connection fails
	}
	defer rabbitConn.Close() // Close RabbitMQ connection

	gmailSender := mail.NewGmailSender(cfg)
	broker := rmqtools.NewMessageBroker(rabbitConn, constants.AmqpExchange, constants.AmqpQueue)
	handlers := handlers.NewMailHandlers(gmailSender, broker)
	broker.ListenForUpdates([]string{constants.RequestEmailVerificationCommand}, handlers.HandleEvent)
}
