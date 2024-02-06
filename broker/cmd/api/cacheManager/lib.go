package cacheManager

import (
	"encoding/json"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
	"time"
)

// userHashKey generates a Redis hash key for user-related data based on the user's Id.
func (cm *cacheManager) userHashKey(uid string) string {
	return fmt.Sprintf("%s-%s", userKey, uid)
}

// userHashKey generates a Redis hash key for user-related data based on the user's Id.
func (cm *cacheManager) accessHKey(uid string) string {
	return fmt.Sprintf("access-%s", uid)
}

// userHashKey generates a Redis hash key for user-related data based on the user's Id.
func (cm *cacheManager) refreshHKey(uid string) string {
	return fmt.Sprintf("refresh-%s", uid)
}

// readCacheTo retrieves the value from the Redis hash and unmarshals it into the provided target.
// It uses the specified key and hash key to read the value from the Redis hash.
func (cm *cacheManager) readCacheTo(readTo interface{}, hashKey, key string) error {
	// Retrieve the value from the Redis hash
	val, err := cm.redisClient.HGet(hashKey, key).Result()
	if err != nil {
		return goErrorHandler.OperationFailure("reading cache", err)
	}

	// Unmarshal the Redis value into the provided readTo
	err = lib.UnmarshalData([]byte(val), readTo)
	if err != nil {
		return err
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
