package lib

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/rmqtools"
	"github.com/go-redis/redis"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

type services struct {
	Redis  *redis.Client
	Rabbit *amqp091.Connection
}

func InitializeServices(cfg config.AppCfg) services {
	// Connect to RabbitMQ server using the provided URL.
	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.RabbitUrl)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// Establish a Redis connection
	redisConn := connectToRedis(cfg.RedisUrl)
	return services{
		Redis:  redisConn,
		Rabbit: rabbitConn,
	}
}
