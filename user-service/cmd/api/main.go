package main

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/server"
	"log"
	"os"
	"time"
)

// main is the entry point of the application.
func main() {
	// Create a background context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure cancellation of the context

	// Load application configuration.
	cfg, err := config.NewConfig()
	if err != nil {
		log.Println(err) // Log error if configuration loading fails
		os.Exit(1)       // Exit program if configuration loading fails
	}

	// Initialize services using the loaded configuration
	services := lib.InitializeServices(ctx, cfg)

	// Close connections when main function exits
	defer services.Rabbit.Close() // Close RabbitMQ connection
	defer services.Redis.Close()  // Close Redis connection
	defer func() {
		// Disconnect from MongoDB and handle error if any
		if err = services.Mongo.Disconnect(ctx); err != nil {
			panic(err) // Panic if MongoDB disconnection fails
		}
	}()

	// Create a new instance of the application using the loaded configuration and service connections, then start it
	server.NewApp(cfg, services.Rabbit, services.Mongo, services.Redis).Start()
}
