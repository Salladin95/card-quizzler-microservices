package main

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/fireBase"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/server"
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

	// Create a background context
	ctx := context.Background()

	// Initialize a Firebase client using the provided configuration
	fireBaseApp := fireBase.NewFireBaseApp(cfg.FireBaseCfg)

	// Connect to the auth
	authClient, err := fireBaseApp.Auth(ctx)
	if err != nil {
		// Log the error and exit the program if connection fails
		log.Fatalf("error getting Auth client: %v\n", err)
		os.Exit(1)
	}

	// Connect to the Firestore database
	firestore, err := fireBaseApp.Firestore(ctx)
	if err != nil {
		// Log the error and exit the program if connection fails
		log.Fatalf("error connecting to Firestore: %v", err)
		os.Exit(1)
	}

	// Ensure the Firestore client is closed when the function completes
	defer firestore.Close()

	// Create a new instance of the application using the loaded configuration and RabbitMQ connection & start it
	server.NewApp(cfg.AppCfg, rabbitConn, firestore, authClient).Start()
}

//consumer, err := rmqtools.NewConsumer(app.rabbit, AmqpExchange, AmqpQueue)
//if err != nil {
//	log.Panic(err)
//}
//err = consumer.Listen([]string{SignInKey, SignUpKey}, handlePayload)
//if err != nil {
//	log.Panic(err)
//}
