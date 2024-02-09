package cacheManager

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"log"
)

func (cm *cacheManager) GetUsers() ([]*entities.UserResponse, error) {
	var cachedUsers []*entities.UserResponse
	// Try to read users from the cache
	err := cm.readCacheByKey(&cachedUsers, cm.userKey)
	if err != nil {
		return nil, err
	}
	// If cache read succeeds, return users from the cache
	log.Println("[api-service] Users retrieved from the cache")
	return cachedUsers, nil
}

func (cm *cacheManager) GetUserById(uid string) (*entities.UserResponse, error) {
	var cachedUser *entities.UserResponse
	// Try to read users from the cache
	err := cm.readCacheByKey(&cachedUser, cm.userHashKey(uid))
	if err != nil {
		return nil, err
	}
	// If cache read succeeds, return users from the cache
	log.Println("[api-service] User retrieved from the cache")
	return cachedUser, nil
}

// GetUserByEmail retrieves a user by their email, either from cache or the underlying repository.
func (cm *cacheManager) GetUserByEmail(email string) (*entities.UserResponse, error) {
	var cachedUser *entities.UserResponse
	// Try to read users from the cache
	err := cm.readCacheByKey(&cachedUser, email)
	if err != nil {
		return nil, err
	}
	// If cache read succeeds, return users from the cache
	log.Println("[api-service] User retrieved from the cache")
	return cachedUser, nil
}
