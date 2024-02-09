package main

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/server"
	"github.com/Salladin95/rmqtools"
	"github.com/go-redis/redis"
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
		log.Println("Error connecting to RabbitMQ:", err)
		os.Exit(1)
	} else {
		log.Println("RabbitMQ connection established successfully!")
	}
	// Ensure the RabbitMQ connection is closed when the main function exits.
	defer rabbitConn.Close()

	// Connect to Mongo database
	mongoClient := connectToMongo(cfg.MongoCfg, ctx)
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Establish a Redis connection
	redisConn := connectToRedis(cfg.AppCfg.RedisUrl)

	// Defer the closure of the Redis connection
	defer redisConn.Close()

	// Create a new instance of the application using the loaded configuration and RabbitMQ connection & start it
	server.NewApp(cfg, rabbitConn, mongoClient, redisConn).Start()
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

// connectToRedis establishes a connection to a Redis server and returns a Redis client.
// It takes the address of the Redis server as a parameter.
func connectToRedis(addr string) *redis.Client {
	// Create a new Redis client with specified options
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		WriteTimeout: 5 * time.Second, // Maximum time to wait for write operations
		ReadTimeout:  5 * time.Second, // Maximum time to wait for read operations
		DialTimeout:  3 * time.Second, // Maximum time to wait for a connection to be established
		MaxRetries:   3,               // Maximum number of retries before giving up on a command
	})
}
