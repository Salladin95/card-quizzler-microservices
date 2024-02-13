package cacheManager

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
)

func (cm *cacheManager) GetUsers(ctx context.Context, uid string) ([]*entities.UserResponse, error) {
	var cachedUsers []*entities.UserResponse
	// Try to read users from the cache
	err := cm.readCacheByKeys(&cachedUsers, cm.userHashKey(uid), usersKey)
	if err != nil {
		return nil, err
	}
	// If cache read succeeds, return users from the cache
	cm.log(ctx, "users has been retrieved from cache", "info", "GetUsers")
	return cachedUsers, nil
}

func (cm *cacheManager) GetUserById(ctx context.Context, uid string) (*entities.UserResponse, error) {
	var cachedUser *entities.UserResponse
	// Try to read users from the cache
	err := cm.readCacheByKeys(&cachedUser, cm.userHashKey(uid), userKey)
	if err != nil {
		return nil, err
	}
	// If cache read succeeds, return users from the cache
	cm.log(ctx, "user has been retrieved from cache", "info", "GetUserById")
	return cachedUser, nil
}
