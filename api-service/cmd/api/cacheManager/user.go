package cacheManager

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
)

func (cm *cacheManager) GetUserById(ctx context.Context, uid string) (*entities.UserResponse, error) {
	var cachedUser *entities.UserResponse
	// Try to read users from the cache
	err := cm.ReadCacheByKeys(&cachedUser, cm.UserHashKey(uid), UserKey)
	if err != nil {
		lib.LogError(err)
		return nil, err
	}
	// If cache read succeeds, return users from the cache
	lib.LogInfo("user has been retrieved from cache")
	return cachedUser, nil
}
