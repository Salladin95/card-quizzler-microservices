package cacheManager

import (
	"encoding/json"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
	"github.com/go-redis/redis"
	"time"
)

// cacheManager represents a manager for handling caching operations using Redis.
type cacheManager struct {
	redisClient *redis.Client  // Redis client for cache operations
	cfg         *config.Config // Application configuration
	userKey     string         // Key for user-related data in the cache
}

// CacheManager is an interface defining methods for caching operations.
type CacheManager interface {
	AccessToken(uid string) (string, error)
	RefreshToken(uid string) (string, error)
	ClearUserData(uid string) error
	SetTokenPair(uid string, tokenPair *entities.TokenPair) error
}

// NewCacheManager creates a new CacheManager instance with the provided Redis client and configuration.
func NewCacheManager(redisClient *redis.Client, cfg *config.Config) CacheManager {
	return &cacheManager{
		redisClient: redisClient,
		cfg:         cfg,
		userKey:     "user", // Default key for user-related data in the cache
	}
}

// SetTokenPair sets access token & refresh token in the cache
// it takes uid & entities.TokenPair as parameters
func (cm *cacheManager) SetTokenPair(uid string, tokenPair *entities.TokenPair) error {
	err := cm.setAccessToken(uid, tokenPair.AccessToken)
	if err != nil {
		return err
	}
	err = cm.setRefreshToken(uid, tokenPair.RefreshToken)
	return err
}

// AccessToken retrieves access token from cache
func (cm *cacheManager) AccessToken(uid string) (string, error) {
	var token string
	err := cm.readCache(&token, cm.userHashKey(uid), cm.accessHashKey(uid))
	if err != nil {
		return "", goErrorHandler.OperationFailure("get access token from cache", err)
	}
	return token, nil
}

// RefreshToken retrieves refresh token from cache
func (cm *cacheManager) RefreshToken(uid string) (string, error) {
	var token string
	err := cm.readCache(&token, cm.userHashKey(uid), cm.refreshHashKey(uid))
	if err != nil {
		return "", goErrorHandler.OperationFailure("get refresh token from cache", err)
	}
	return token, nil
}

// ClearUserData drops user related cache
func (cm *cacheManager) ClearUserData(uid string) error {
	return cm.redisClient.Del(cm.userHashKey(uid)).Err()
}

// SetAccessToken sets the access token for a user in the cache.
func (cm *cacheManager) setAccessToken(uid, token string) error {
	err := cm.setCacheInPipeline(cm.userHashKey(uid), cm.accessHashKey(uid), token, cm.cfg.JwtCfg.AccessTokenExpTime)
	if err != nil {
		return goErrorHandler.OperationFailure("set access token cache", err)
	}
	return nil
}

// SetRefreshToken sets the refresh token for a user in the cache.
func (cm *cacheManager) setRefreshToken(uid, token string) error {
	err := cm.setCacheInPipeline(cm.userHashKey(uid), cm.refreshHashKey(uid), token, cm.cfg.JwtCfg.AccessTokenExpTime)
	if err != nil {
		return goErrorHandler.OperationFailure("set access token cache", err)
	}
	return nil
}

// readCache retrieves the value from the Redis hash and unmarshals it into the provided target.
// It uses the specified key and hash key to read the value from the Redis hash.
func (cm *cacheManager) readCache(target interface{}, key, hashKey string) error {
	// Retrieve the value from the Redis hash
	val, err := cm.redisClient.HGet(key, hashKey).Result()
	if err != nil {
		return goErrorHandler.OperationFailure("reading cache", err)
	}

	// Unmarshal the Redis value into the provided target
	err = lib.UnmarshalData([]byte(val), target)
	if err != nil {
		return goErrorHandler.OperationFailure("unmarshal cache value", err)
	}

	return nil
}

// setCacheInPipeline sets data in the cache using a Redis pipeline to perform multiple operations in a single round trip.
func (cm *cacheManager) setCacheInPipeline(key string, hash string, data any, exp time.Duration) error {
	pipe := cm.redisClient.Pipeline()
	defer pipe.Close()

	data, err := json.Marshal(data)

	if err != nil {
		return goErrorHandler.OperationFailure("marshal data before setting cache", err)
	}

	// Set hash field
	pipe.HSet(key, hash, data)

	// Set expiration time
	pipe.Expire(key, exp)

	// Execute pipeline
	_, err = pipe.Exec()

	if err != nil {
		return goErrorHandler.OperationFailure("set cache", err)
	}
	return nil
}

// userHashKey generates a Redis hash key for user-related data based on the user's Id.
func (cm *cacheManager) userHashKey(uid string) string {
	return fmt.Sprintf("%s-%s", cm.userKey, uid)
}

// accessHashKey generates a Redis hash key for the access token based on the user's Id.
func (cm *cacheManager) accessHashKey(uid string) string {
	return fmt.Sprintf("access-token-%s", uid)
}

// refreshHashKey generates a Redis hash key for the refresh token based on the user's Id.
func (cm *cacheManager) refreshHashKey(uid string) string {
	return fmt.Sprintf("refresh-token-%s", uid)
}
