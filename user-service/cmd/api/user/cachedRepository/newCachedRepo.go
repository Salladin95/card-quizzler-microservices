package cachedRepository

import (
	userRepo "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/repository"
	"github.com/go-redis/redis"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

func NewCachedUserRepo(
	redisClient *redis.Client,
	rabbitConn *amqp091.Connection,
	userRep userRepo.Repository,
) CachedRepository {
	return &cachedRepository{
		redisClient: redisClient,
		repo:        userRep,
		userKey:     "user",
		exp:         60 * time.Minute,
		rabbitConn:  rabbitConn,
	}
}
