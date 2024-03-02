package cacheManager

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
)

// ReadCacheKeys retrieves the value from the Redis hash and unmarshals it into the provided readTo parameter.
// It uses the specified key and hash key to read the value from the Redis hash.
// Note: The readTo parameter must be a pointer.
func (cm *cacheManager) ReadCacheKeys(readTo interface{}, key, hashKey string) error {
	// Retrieve the value from the Redis hash
	val, err := cm.redisClient.HGet(key, hashKey).Result()
	if err != nil {
		return goErrorHandler.OperationFailure("read cache", err)
	}

	// Unmarshal the Redis value into the provided readTo
	err = lib.UnmarshalData([]byte(val), readTo)
	if err != nil {
		return err
	}
	return nil
}

// ReadCacheByKey retrieves the value from Redis using a key and returns the result.
// Note: The readTo parameter must be a pointer.
func (cm *cacheManager) ReadCacheByKey(readTo interface{}, key string) error {
	// Retrieve the value from the Redis hash
	val, err := cm.redisClient.Get(key).Result()
	if err != nil {
		return goErrorHandler.OperationFailure("read cache", err)
	}

	// Unmarshal the Redis value into the provided readTo
	err = lib.UnmarshalData([]byte(val), readTo)
	if err != nil {
		return err
	}
	return nil
}

func (cm *cacheManager) SetCacheByKey(key string, data []byte) error {
	// Set the marshalled data in the Redis cache with the specified expiration time
	err := cm.redisClient.Set(key, data, cm.exp).Err()
	if err != nil {
		return goErrorHandler.OperationFailure(fmt.Sprintf("set cache by key - %s", key), err)
	}
	return nil
}

func (cm *cacheManager) SetCacheByKeys(key string, hash string, data []byte) error {
	// Create a new Redis pipeline
	pipe := cm.redisClient.Pipeline()
	defer pipe.Close()

	// Set hash field in the Redis cache
	pipe.HSet(key, hash, data)

	// Set expiration time for the cache
	pipe.Expire(key, cm.exp)

	// Execute the pipeline to perform multiple operations in a single round trip
	if _, err := pipe.Exec(); err != nil {
		return goErrorHandler.OperationFailure("set cache", err)
	}

	return nil
}

func (cm *cacheManager) ClearCacheByKey(key string) error {
	// Delete the cache entry using the Redis client
	return cm.redisClient.Del(key).Err()
}

func (cm *cacheManager) ClearCacheByKeys(key1, key2 string) error {
	// Delete the cache entry using the Redis client
	return cm.redisClient.HDel(key1, key2).Err()
}
