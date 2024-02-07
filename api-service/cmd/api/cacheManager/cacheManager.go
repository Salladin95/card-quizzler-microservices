package cacheManager

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/goErrorHandler"
	"github.com/go-redis/redis"
	"github.com/labstack/gommon/log"
)

const (
	userKey = "user"
)

// cacheManager represents a manager for handling caching operations using Redis.
type cacheManager struct {
	redisClient *redis.Client  // Redis client for cache operations
	cfg         *config.Config // Application configuration
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
	log.Infof("token - %s has been set to cache", token)
	return nil
}

// SetRefreshToken sets the refresh token for a user in the cache.
func (cm *cacheManager) setRefreshToken(uid, token string) error {
	err := cm.redisClient.Set(cm.refreshHKey(uid), token, cm.cfg.JwtCfg.AccessTokenExpTime).Err()
	if err != nil {
		return goErrorHandler.OperationFailure("set access token cache", err)
	}
	log.Infof("token - %s has been set to cache", token)
	return nil
}

// AccessToken retrieves access token from cache
func (cm *cacheManager) AccessToken(uid string) (string, error) {
	token, err := cm.redisClient.Get(cm.accessHKey(uid)).Result()
	if err != nil {
		return "", goErrorHandler.OperationFailure("get access token from cache", err)
	}
	log.Infof("token - %s has been extracted from cache", token)
	return token, nil
}

// RefreshToken retrieves refresh token from cache
func (cm *cacheManager) RefreshToken(uid string) (string, error) {
	token, err := cm.redisClient.Get(cm.refreshHKey(uid)).Result()
	if err != nil {
		return "", goErrorHandler.OperationFailure("get refresh token from cache", err)
	}
	log.Infof("token - %s has been extracted from cache", token)
	return token, nil
}

// ClearUserData drops user related cache
func (cm *cacheManager) ClearUserData(uid string) error {
	return cm.redisClient.Del(cm.userHashKey(uid)).Err()
}
