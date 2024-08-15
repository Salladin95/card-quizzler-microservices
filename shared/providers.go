package lib

import (
	"github.com/go-redis/redis"
	"time"
)

// ConnectToRedis establishes a connection to a Redis server and returns a Redis client.
// It takes the address of the Redis server as a parameter.
func ConnectToRedis(addr string) *redis.Client {
	// Create a new Redis client with specified options
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		WriteTimeout: 5 * time.Second, // Maximum time to wait for write operations
		ReadTimeout:  5 * time.Second, // Maximum time to wait for read operations
		DialTimeout:  3 * time.Second, // Maximum time to wait for a connection to be established
		MaxRetries:   3,               // Maximum number of retries before giving up on a command
	})
}
