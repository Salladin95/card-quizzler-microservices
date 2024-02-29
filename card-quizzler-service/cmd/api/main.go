package main

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/server"
	"github.com/Salladin95/rmqtools"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func main() {
	s := echo.New()
	cfg, err := config.NewConfig()

	if err != nil {
		s.Logger.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Connect to RabbitMQ s using the provided URL.
	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.RabbitUrl)
	if err != nil {
		log.Println(err) // Log error if RabbitMQ connection fails
		os.Exit(1)       // Exit program if RabbitMQ connection fails
	}

	redis := connectToRedis(cfg.RedisUrl)

	server.NewApp(cfg, rabbitConn, db, redis).Start()
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
