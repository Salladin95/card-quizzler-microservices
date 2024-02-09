package cacheManager

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/goErrorHandler"
	"github.com/go-redis/redis"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

const (
	userKey = "user"
)

// cacheManager represents a manager for handling caching operations using Redis.
type cacheManager struct {
	redisClient *redis.Client       // Redis client for cache operations
	cfg         *config.Config      // Application configuration
	rabbitConn  *amqp091.Connection // rabbitConn is the AMQP connection used for firing events.
	exp         time.Duration       // exp is the expiration time for cached data.
	userKey     string              // userKey is the key used to store and retrieve user data from the cache.
}

// CacheManager is an interface defining methods for caching operations.
type CacheManager interface {
	AccessToken(uid string) (string, error)
	RefreshToken(uid string) (string, error)
	ClearDataByUID(uid string) error
	SetTokenPair(uid string, tokenPair *entities.TokenPair) error
	GetUsers() ([]*entities.UserResponse, error)
	GetUserById(uid string) (*entities.UserResponse, error)
	GetUserByEmail(email string) (*entities.UserResponse, error)
	ListenForUpdates()
}

// NewCacheManager creates a new CacheManager instance with the provided Redis client and configuration.
func NewCacheManager(redisClient *redis.Client, cfg *config.Config, rabbitConn *amqp091.Connection) CacheManager {
	return &cacheManager{
		redisClient: redisClient,
		cfg:         cfg,
		exp:         60 * time.Minute,
		rabbitConn:  rabbitConn,
		userKey:     "api-service_user",
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

// SetAccessToken sets the access token for a user in the cache.
func (cm *cacheManager) setAccessToken(uid, token string) error {
	err := cm.redisClient.Set(cm.accessHKey(uid), token, cm.cfg.JwtCfg.AccessTokenExpTime).Err()
	if err != nil {
		return goErrorHandler.OperationFailure("set access token cache", err)
	}
	return nil
}

// SetRefreshToken sets the refresh token for a user in the cache.
func (cm *cacheManager) setRefreshToken(uid, token string) error {
	err := cm.redisClient.Set(cm.refreshHKey(uid), token, cm.cfg.JwtCfg.AccessTokenExpTime).Err()
	if err != nil {
		return goErrorHandler.OperationFailure("set access token cache", err)
	}
	return nil
}

// AccessToken retrieves access token from cache
func (cm *cacheManager) AccessToken(uid string) (string, error) {
	token, err := cm.redisClient.Get(cm.accessHKey(uid)).Result()
	if err != nil {
		return "", goErrorHandler.OperationFailure("get access token from cache", err)
	}
	return token, nil
}

// RefreshToken retrieves refresh token from cache
func (cm *cacheManager) RefreshToken(uid string) (string, error) {
	token, err := cm.redisClient.Get(cm.refreshHKey(uid)).Result()
	if err != nil {
		return "", goErrorHandler.OperationFailure("get refresh token from cache", err)
	}
	return token, nil
}

// ClearDataByUID ClearUserData drops user related cache
func (cm *cacheManager) ClearDataByUID(uid string) error {
	return cm.redisClient.Del(cm.userHashKey(uid)).Err()
}

// ClearDataByKey drops data by provided key
func (cm *cacheManager) ClearDataByKey(key string) error {
	return cm.redisClient.Del(key).Err()
}
