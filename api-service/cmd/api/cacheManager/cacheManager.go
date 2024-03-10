package cacheManager

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
	"github.com/Salladin95/rmqtools"
	"github.com/go-redis/redis"
	"time"
)

// cacheManager represents a manager for handling caching operations using Redis.
type cacheManager struct {
	redisClient *redis.Client          // Redis client for cache operations
	cfg         *config.Config         // Application configuration
	broker      rmqtools.MessageBroker // MessageBroker instance
	exp         time.Duration          // exp is the expiration time for cached data.
}

// CacheManager is an interface defining methods for caching operations.
type CacheManager interface {
	AccessToken(uid string) (string, error)
	RefreshToken(uid string) (string, error)
	ClearUserRelatedCache(uid string) error
	ClearCacheByKeys(key1, key2 string) error
	ClearCacheByKey(key string) error
	SetTokenPair(uid string, tokenPair *entities.TokenPair) error
	GetUserById(ctx context.Context, uid string) (*entities.UserResponse, error)
	SetCacheByKeys(key string, hash string, data []byte, exp time.Duration) error
	SetCacheByKey(key string, data []byte) error
	ReadCacheByKey(readTo interface{}, key string) error
	ReadCacheByKeys(readTo interface{}, key, hashKey string) error
	UserHashKey(uid string) string
	Exp() time.Duration
}

const (
	RootKey          = "hash:api-service_user"
	TokensKey        = "hash:token-pair"
	UsersKey         = "hash:users"
	UserKey          = "hash:user"
	Folders          = "hash:folders"
	Folder           = "hash:folder"
	Modules          = "hash:modules"
	Module           = "hash:module"
	DifficultModules = "difficult:modules"
)

func FolderKey(id string) string {
	return fmt.Sprintf("folder:%s", id)
}

func ModuleKey(id string) string {
	return fmt.Sprintf("module:%s", id)
}

func FoldersKey(uid string) string {
	return fmt.Sprintf("%s:%s", Folders, uid)
}

func ModulesKey(uid string) string {
	return fmt.Sprintf("%s:%s", Modules, uid)
}

// NewCacheManager creates a new CacheManager instance with the provided Redis client and configuration.
func NewCacheManager(
	redisClient *redis.Client,
	cfg *config.Config,
	broker rmqtools.MessageBroker,
) CacheManager {
	return &cacheManager{
		redisClient: redisClient,
		cfg:         cfg,
		exp:         60 * time.Minute,
		broker:      broker,
	}
}

func (cm *cacheManager) Exp() time.Duration {
	return cm.exp
}

// ClearUserRelatedCache drops user related cache
func (cm *cacheManager) ClearUserRelatedCache(uid string) error {
	return cm.redisClient.Del(cm.UserHashKey(uid)).Err()
}

// ClearCacheByKeys drops specified cache
func (cm *cacheManager) ClearCacheByKeys(key string, key2 string) error {
	if err := cm.redisClient.HDel(key, key2).Err(); err != nil {
		cm.log(context.Background(), err.Error(), "error", "ClearCacheByKeys")
		return goErrorHandler.OperationFailure("clear cache", err)
	}
	return nil
}

// ClearCacheByKey drops specified cache
func (cm *cacheManager) ClearCacheByKey(key string) error {
	if err := cm.redisClient.Del(key).Err(); err != nil {
		cm.log(context.Background(), err.Error(), "error", "ClearCacheByKeys")
		return goErrorHandler.OperationFailure("clear cache", err)
	}
	return nil
}

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

// SetCacheByKeys sets data in the cache using a Redis pipeline to perform multiple operations in a single round trip.
// It takes the specified key, hash, data and exp as parameters, marshals the data into JSON format,
// It returns an error if any issues occur during the marshaling or cache setting process.
func (cm *cacheManager) SetCacheByKeys(key string, hash string, data []byte, exp time.Duration) error {
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

// log sends a log message to the message broker.
func (cm *cacheManager) log(ctx context.Context, message, level, method string) {
	var log entities.LogMessage

	// Push log message to the message broker
	cm.broker.PushToQueue(
		ctx,
		constants.LogCommand, // Specify the log command constant
		// Generate log message with provided details
		log.GenerateLog(message, level, method, "cache manager"),
	)
}
