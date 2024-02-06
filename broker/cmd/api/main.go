package main

import (
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/server"
	"github.com/Salladin95/rmqtools"
	"github.com/go-redis/redis"
	"log"
	"os"
	"time"
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

	// Establish a Redis connection
	redisConn := connectToRedis(cfg.AppCfg.RedisUrl)

	// Defer the closure of the Redis connection
	defer redisConn.Close()

	// Create a new instance of the application using the loaded configuration and RabbitMQ connection & start it
	server.NewApp(cfg, rabbitConn, redisConn).Start()
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
