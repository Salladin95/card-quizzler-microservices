package cacheManager

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/messageBroker"
	"github.com/go-redis/redis"
	"time"
)

const (
	userKey = "user"
)

// cacheManager represents a manager for handling caching operations using Redis.
type cacheManager struct {
	redisClient   *redis.Client               // Redis client for cache operations
	cfg           *config.Config              // Application configuration
	messageBroker messageBroker.MessageBroker // messageBroker instance
	exp           time.Duration               // exp is the expiration time for cached data.
	userKey       string                      // userKey is the key used to store and retrieve user data from the cache.
}

// CacheManager is an interface defining methods for caching operations.
type CacheManager interface {
	AccessToken(uid string) (string, error)
	RefreshToken(uid string) (string, error)
	ClearDataByUID(uid string) error
	SetTokenPair(uid string, tokenPair *entities.TokenPair) error
	GetUsers(ctx context.Context) ([]*entities.UserResponse, error)
	GetUserById(ctx context.Context, uid string) (*entities.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.UserResponse, error)
	ListenForUpdates()
}

// NewCacheManager creates a new CacheManager instance with the provided Redis client and configuration.
func NewCacheManager(
	redisClient *redis.Client,
	cfg *config.Config,
	broker messageBroker.MessageBroker,
) CacheManager {
	return &cacheManager{
		redisClient:   redisClient,
		cfg:           cfg,
		exp:           60 * time.Minute,
		messageBroker: broker,
		userKey:       "api-service_user",
	}
}

// ClearDataByUID ClearUserData drops user related cache
func (cm *cacheManager) ClearDataByUID(uid string) error {
	return cm.redisClient.Del(cm.userHashKey(uid)).Err()
}

// ClearDataByKey drops data by provided key
func (cm *cacheManager) ClearDataByKey(key string) error {
	return cm.redisClient.Del(key).Err()
}
