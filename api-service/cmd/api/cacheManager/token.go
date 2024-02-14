package cacheManager

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
)

// AccessToken retrieves access token from cache
func (cm *cacheManager) AccessToken(uid string) (string, error) {
	var tokenPair entities.TokenPair
	err := cm.ReadCacheByKeys(&tokenPair, cm.UserHashKey(uid), TokensKey)
	if err != nil {
		return "", goErrorHandler.OperationFailure("get access token from cache", err)
	}
	return tokenPair.AccessToken, nil
}

// RefreshToken retrieves refresh token from cache
func (cm *cacheManager) RefreshToken(uid string) (string, error) {
	var tokenPair entities.TokenPair
	err := cm.ReadCacheByKeys(&tokenPair, cm.UserHashKey(uid), TokensKey)
	if err != nil {
		return "", goErrorHandler.OperationFailure("get refresh token from cache", err)
	}
	return tokenPair.RefreshToken, nil
}

// SetTokenPair sets access token & refresh token in the cache
// it takes uid & entities.TokenPair as parameters
func (cm *cacheManager) SetTokenPair(uid string, tokenPair *entities.TokenPair) error {
	// set refresh token exp time, because it's lives longer than access token exp time
	parsedData, err := lib.MarshalData(tokenPair)
	if err != nil {
		return err
	}
	err = cm.SetCacheInPipeline(cm.UserHashKey(uid), TokensKey, parsedData, cm.cfg.JwtCfg.RefreshTokenExpTime)
	return err
}
