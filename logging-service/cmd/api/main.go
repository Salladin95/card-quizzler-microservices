package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/models"
	"github.com/Salladin95/rmqtools"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func main() {
	// Create a background context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load application configuration.
	cfg, err := config.NewConfig()

	// Connect to RabbitMQ server using the provided URL.
	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.RabbitUrl)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Connect to Mongo database
	mongoClient := connectToMongo(cfg.MongoCfg, ctx)
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	consumer, err := rmqtools.NewConsumer(
		rabbitConn,
		constants.AmqpExchange,
		constants.AmqpQueue,
	)
	if err != nil {
		// TODO: REPLACE
		fmt.Println(err)
	}
	err = consumer.Listen(
		[]string{
			constants.LogCommand,
		},
		func(routingKey string, payload []byte) {
			var logMessage entities.LogMessage
			err := json.Unmarshal(payload, &logMessage)
			if err != nil {
				fmt.Println("failed to unmarshall logMessage")
				return
			}

			err = logMessage.Verify()
			if err != nil {
				fmt.Printf("log-message has failed validation - %v", err)
				return
			}

			fmt.Printf(
				"[logging service][routingkey - %s][message from - %s][method - %s][level - %s][description - %s] message - %s\n",
				routingKey, logMessage.FromService, logMessage.Method, logMessage.Level, logMessage.Name, logMessage.Message,
			)
			// Create a background context
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			models := models.NewModels(mongoClient)

			err = models.LogEntry.Insert(ctx, logMessage)
			if err != nil {
				fmt.Printf("failed to insert log - %v", err)
			}
		},
	)
	if err != nil {
		// TODO: REPLACE
		fmt.Println(err)
	}
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
