package lib

import (
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

func InitializeServices(cfg ServicesCfg) services {
	// Connect to RabbitMQ server using the provided URL.
	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.RabbitUrl)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// Establish a Redis connection
	redisConn := ConnectToRedis(cfg.RedisUrl)
	return services{
		Redis:  redisConn,
		Rabbit: rabbitConn,
	}
}
