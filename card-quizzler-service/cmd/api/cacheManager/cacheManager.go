package cacheManager

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type CacheManager interface {
	ReadCacheKeys(readTo interface{}, key, hashKey string) error
	ReadCacheByKey(readTo interface{}, key string) error
	SetCacheByKeys(key1, key2 string, data []byte) error
	SetCacheByKey(key string, data []byte) error
	ClearCacheByKeys(key1, key2 string) error
	ClearCacheByKey(key string) error
	RootKey(uid string) string
}

type cacheManager struct {
	redisClient *redis.Client
	exp         time.Duration
}

func (cm *cacheManager) RootKey(uid string) string {
	return fmt.Sprintf("hash:card-quiz-cache-%s", uid)
}

func NewCacheManager(redis *redis.Client, exp time.Duration) CacheManager {
	return &cacheManager{
		redisClient: redis,
		exp:         exp,
	}
}
