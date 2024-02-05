package main

import (
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/server"
	"github.com/Salladin95/rmqtools"
	"log"
	"os"
)

// main is the entry point of the application.
func main() {
	// Load application configuration.
	cfg, err := config.NewConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Connect to RabbitMQ server using the provided URL.
	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.AppCfg.RabbitUrl)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Ensure the RabbitMQ connection is closed when the main function exits.
	defer rabbitConn.Close()

	// Create a new instance of the application using the loaded configuration and RabbitMQ connection & start it
	server.NewApp(cfg, rabbitConn).Start()
}
