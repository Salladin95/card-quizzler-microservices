package main

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/server"
	"time"
)

func main() {
	// Create a background context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load application configuration.
	cfg, err := config.NewConfig()

	// Initialize services using the loaded configuration
	services := lib.InitializeServices(ctx, cfg)

	// Close connections when main function exits
	defer services.Rabbit.Close() // Close RabbitMQ connection
	defer func() {
		// Disconnect from MongoDB and handle error if any
		if err = services.Mongo.Disconnect(ctx); err != nil {
			panic(err) // Panic if MongoDB disconnection fails
		}
	}()

	server.NewApp(services.Rabbit, services.Mongo, cfg).Start()
}
