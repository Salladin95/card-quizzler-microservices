package cachedRepository

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
)

// readCacheByHashedKey retrieves the value from the Redis hash and unmarshals it into the provided readTo parameter.
// It uses the specified key and hash key to read the value from the Redis hash.
// Note: The readTo parameter must be a pointer.
func (cr *cachedRepository) readCacheByHashedKey(readTo interface{}, key, hashKey string) error {
	// Retrieve the value from the Redis hash
	val, err := cr.redisClient.HGet(key, hashKey).Result()
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

// readCacheByKey retrieves the value from Redis using a key and returns the result.
// Note: The readTo parameter must be a pointer.
func (cr *cachedRepository) readCacheByKey(readTo interface{}, key string) error {
	// Retrieve the value from the Redis hash
	val, err := cr.redisClient.Get(key).Result()
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

// setCacheByKey sets data in the Redis cache for the specified key.
// It marshals the data into JSON format and stores it in the Redis hash using the key.
// The expiration time for the cache is determined by the configured expiration duration in the cached repository.
// It returns an error if any issues occur during the marshaling or cache setting process.
func (cr *cachedRepository) setCacheByKey(key string, data interface{}) error {
	// Marshal the data into JSON format
	marshalledData, err := lib.MarshalData(data)
	if err != nil {
		return err
	}

	// Set the marshalled data in the Redis cache with the specified expiration time
	err = cr.redisClient.Set(key, marshalledData, cr.exp).Err()
	if err != nil {
		return goErrorHandler.OperationFailure(fmt.Sprintf("set cache by key - %s", key), err)
	}
	return nil
}

// setCacheInPipeline sets data in the cache using a Redis pipeline to perform multiple operations in a single round trip.
// It takes the specified key, hash, and data as parameters, marshals the data into JSON format,
// and uses a Redis pipeline to set the hash field and expiration time in the cache.
// It returns an error if any issues occur during the marshaling or cache setting process.
func (cr *cachedRepository) setCacheInPipeline(key string, hash string, data interface{}) error {
	// Create a new Redis pipeline
	pipe := cr.redisClient.Pipeline()
	defer pipe.Close()

	// Marshal the data into JSON format
	marshalledData, err := lib.MarshalData(data)
	if err != nil {
		return err
	}

	// Set hash field in the Redis cache
	pipe.HSet(key, hash, marshalledData)

	// Set expiration time for the cache
	pipe.Expire(key, cr.exp)

	// Execute the pipeline to perform multiple operations in a single round trip
	_, err = pipe.Exec()
	if err != nil {
		return goErrorHandler.OperationFailure("set cache", err)
	}

	return nil
}

func (cr *cachedRepository) userHashKey(key string) string {
	return fmt.Sprintf("hash:%s", key)
}

// clearCacheByKey drops the cache associated with the given key.
func (cr *cachedRepository) clearCacheByKey(key string) error {
	// Delete the cache entry using the Redis client
	return cr.redisClient.Del(key).Err()
}

// clearCacheByKey drops the cache associated with the given key.
func (cr *cachedRepository) clearCacheByKeys(key1, key2 string) error {
	// Delete the cache entry using the Redis client
	return cr.redisClient.HDel(key1, key2).Err()
}
