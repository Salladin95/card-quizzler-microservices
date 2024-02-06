package cachedRepository

import (
	userRepo "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/repository"
	"github.com/go-redis/redis"
	"time"
)

func NewCachedUserRepo(redisClient *redis.Client, userRep userRepo.Repository) CachedRepository {
	return &cachedRepository{
		redisClient: redisClient,
		repo:        userRep,
		userKey:     "user",
		exp:         60 * time.Minute,
	}
}
