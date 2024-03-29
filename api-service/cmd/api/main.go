package main

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/server"
	"github.com/Salladin95/rmqtools"
	"log/slog"
	"os"
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

// main is the entry point of the application.
func main() {
	// Load application configuration.
	cfg, err := config.NewConfig()
	if err != nil {
		lib.LogError(err)
		os.Exit(1)
	}

	services := lib.InitializeServices(cfg.AppCfg)
	// Ensure the RabbitMQ connection is closed when the main function exits.
	defer services.Rabbit.Close()
	// Defer the closure of the Redis connection
	defer services.Redis.Close()

	broker := rmqtools.NewMessageBroker(services.Rabbit, constants.AmqpExchange, constants.AmqpQueue)

	// Create a new instance of the application using the loaded configuration and RabbitMQ connection & start it
	server.NewApp(cfg, services.Redis, broker).Start()
}
