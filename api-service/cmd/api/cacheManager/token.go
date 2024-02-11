package cacheManager

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/goErrorHandler"
)

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
