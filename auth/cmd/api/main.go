package main

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/server"
	"github.com/Salladin95/rmqtools"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

// main is the entry point of the application.
func main() {
	// Create a background context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	mongoClient := connectToMongo(cfg.MongoCfg, ctx)
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Create a new instance of the application using the loaded configuration and RabbitMQ connection & start it
	server.NewApp(cfg, rabbitConn, mongoClient).Start()
}

func connectToMongo(mongoCfg config.MongoCfg, ctx context.Context) *mongo.Client {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoCfg.MongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: mongoCfg.MongoUsername,
		Password: mongoCfg.MongoUserPassword,
	})

	// connect
	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("failed connect to mongo - %v\n", err)
		os.Exit(1)
	}

	log.Println("Connected to mongo!")

	return c
}

//consumer, err := rmqtools.NewConsumer(app.rabbit, AmqpExchange, AmqpQueue)
//if err != nil {
//	log.Panic(err)
//}
//err = consumer.Listen([]string{SignInKey, SignUpKey}, handlePayload)
//if err != nil {
//	log.Panic(err)
//}
