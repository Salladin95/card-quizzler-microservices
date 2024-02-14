package cacheManager

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/rmqtools"
	"github.com/go-redis/redis"
	"time"
)

// cacheManager represents a manager for handling caching operations using Redis.
type cacheManager struct {
	redisClient *redis.Client          // Redis client for cache operations
	cfg         *config.Config         // Application configuration
	broker      rmqtools.MessageBroker // MessageBroker instance
	exp         time.Duration          // exp is the expiration time for cached data.
}

// CacheManager is an interface defining methods for caching operations.
type CacheManager interface {
	AccessToken(uid string) (string, error)
	RefreshToken(uid string) (string, error)
	ClearUserRelatedCache(uid string) error
	ClearCacheByKeys(key1, key2 string) error
	SetTokenPair(uid string, tokenPair *entities.TokenPair) error
	GetUserById(ctx context.Context, uid string) (*entities.UserResponse, error)
	SetCacheInPipeline(key string, hash string, data []byte, exp time.Duration) error
	SetCacheByKey(key string, data []byte) error
	ReadCacheByKey(readTo interface{}, key string) error
	ReadCacheByKeys(readTo interface{}, key, hashKey string) error
	UserHashKey(uid string) string
	Exp() time.Duration
}

// NewCacheManager creates a new CacheManager instance with the provided Redis client and configuration.
func NewCacheManager(
	redisClient *redis.Client,
	cfg *config.Config,
	broker rmqtools.MessageBroker,
) CacheManager {
	return &cacheManager{
		redisClient: redisClient,
		cfg:         cfg,
		exp:         60 * time.Minute,
		broker:      broker,
	}
}

func (cm *cacheManager) Exp() time.Duration {
	return cm.exp
}

// ClearUserRelatedCache drops user related cache
func (cm *cacheManager) ClearUserRelatedCache(uid string) error {
	return cm.redisClient.Del(cm.UserHashKey(uid)).Err()
}

// ClearCacheByKeys drops specified cache
func (cm *cacheManager) ClearCacheByKeys(key string, key2 string) error {
	return cm.redisClient.HDel(key, key2).Err()
}

// log sends a log message to the message broker.
func (cm *cacheManager) log(ctx context.Context, message, level, method string) {
	var log entities.LogMessage

	// Push log message to the message broker
	cm.broker.PushToQueue(
		ctx,
		constants.LogCommand, // Specify the log command constant
		// Generate log message with provided details
		log.GenerateLog(message, level, method, "cache manager"),
	)
}
