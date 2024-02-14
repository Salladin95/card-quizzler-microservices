package cacheManager

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
	"time"
)

// UserHashKey generates a Redis hash key for user-related data based on the user's Id.
func (cm *cacheManager) UserHashKey(uid string) string {
	return fmt.Sprintf("%s-%s", RootKey, uid)
}

// ReadCacheByKeys retrieves the value from the Redis hash and unmarshals it into the provided readTo parameter.
// It uses the specified key and hash key to read the value from the Redis hash.
// Note: The readTo parameter must be a pointer.
func (cm *cacheManager) ReadCacheByKeys(readTo interface{}, key, hashKey string) error {
	// Retrieve the value from the Redis hash
	val, err := cm.redisClient.HGet(key, hashKey).Result()
	if err != nil {
		return fmt.Errorf("read cache - %v", err)
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

// SetCacheByKey sets data in the Redis cache for the specified key.
// It marshals the data into JSON format and stores it in the Redis hash using the key.
// The expiration time for the cache is determined by the configured expiration duration in the cached repository.
// It returns an error if any issues occur during the marshaling or cache setting process.
func (cm *cacheManager) SetCacheByKey(key string, data []byte) error {
	// Set the marshalled data in the Redis cache with the specified expiration time
	err := cm.redisClient.Set(key, data, cm.exp).Err()
	if err != nil {
		return goErrorHandler.OperationFailure(fmt.Sprintf("set cache by key - %s", key), err)
	}
	return nil
}

// SetCacheInPipeline sets data in the cache using a Redis pipeline to perform multiple operations in a single round trip.
// It takes the specified key, hash, data and exp as parameters, marshals the data into JSON format,
// It returns an error if any issues occur during the marshaling or cache setting process.
func (cm *cacheManager) SetCacheInPipeline(key string, hash string, data []byte, exp time.Duration) error {
	// Create a new Redis pipeline
	pipe := cm.redisClient.Pipeline()
	defer pipe.Close()

	// Set hash field in the Redis cache
	pipe.HSet(key, hash, data)

	// Set expiration time for the cache
	pipe.Expire(key, exp)

	// Execute the pipeline to perform multiple operations in a single round trip
	_, err := pipe.Exec()
	if err != nil {
		return goErrorHandler.OperationFailure("set cache", err)
	}

	return nil
}
