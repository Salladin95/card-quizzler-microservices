package cacheManager

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
	"github.com/Salladin95/rmqtools"
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

// readCacheByHashedKey retrieves the value from the Redis hash and unmarshals it into the provided readTo parameter.
// It uses the specified key and hash key to read the value from the Redis hash.
// Note: The readTo parameter must be a pointer.
func (cm *cacheManager) readCacheByHashedKey(readTo interface{}, key, hashKey string) error {
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

// readCacheByKey retrieves the value from Redis using a key and returns the result.
// Note: The readTo parameter must be a pointer.
func (cm *cacheManager) readCacheByKey(readTo interface{}, key string) error {
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

// setCacheByHashKeyInPipeline sets data in the cache using a Redis pipeline to perform multiple operations in a single round trip.
// It takes the specified key, hash, and data as parameters, marshals the data into JSON format,
// and uses a Redis pipeline to set the hash field and expiration time in the cache.
// It returns an error if any issues occur during the marshaling or cache setting process.
func (cm *cacheManager) setCacheByHashKeyInPipeline(key string, hash string, data interface{}) error {
	// Create a new Redis pipeline
	pipe := cm.redisClient.Pipeline()
	defer pipe.Close()

	// Marshal the data into JSON format
	marshalledData, err := lib.MarshalData(data)
	if err != nil {
		return err
	}

	// Set hash field in the Redis cache
	pipe.HSet(key, hash, marshalledData)

	// Set expiration time for the cache
	pipe.Expire(key, cm.exp)

	// Execute the pipeline to perform multiple operations in a single round trip
	_, err = pipe.Exec()
	if err != nil {
		return goErrorHandler.OperationFailure("set cache", err)
	}

	return nil
}

// pushToQueue pushes data to a named queue in RabbitMQ using an EventEmitter.
// It takes a context, the name of the queue, and the data to be pushed.
// It returns an error if any operation fails.
func (cm *cacheManager) pushToQueue(ctx context.Context, name string, data any) error {
	// Create a new EventEmitter for the specified AMQP exchange.
	emitter, err := rmqtools.NewEventEmitter(cm.rabbitConn, constants.AmqpExchange)
	if err != nil {
		return goErrorHandler.OperationFailure("create event emitter", err)
	}

	// Marshal the data into JSON format.
	mData, err := lib.MarshalData(data)
	if err != nil {
		return err
	}

	// Push the data to the named queue using the EventEmitter.
	err = emitter.Push(ctx, name, mData)
	if err != nil {
		return goErrorHandler.OperationFailure("push event", err)
	}

	// Print a message indicating that the event has been pushed successfully.
	fmt.Printf("event - %s has been pushed\n", name)

	return nil
}