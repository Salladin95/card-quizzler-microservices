package lib

import (
	"context"
	"encoding/json"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/config"
	"github.com/Salladin95/goErrorHandler"
	"github.com/Salladin95/rmqtools"
	"github.com/go-redis/redis"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

// UnmarshalData unmarshals JSON data into the provided unmarshalTo interface.
// It returns an error if any issues occur during the unmarshaling process.
// Not - unmarshalTo must be pointer !!!
func UnmarshalData(data []byte, unmarshalTo interface{}) error {
	err := json.Unmarshal(data, unmarshalTo)
	if err != nil {
		return goErrorHandler.OperationFailure("unmarshal data", err)
	}
	return nil
}

// MarshalData marshals data into a JSON-encoded byte slice.
// It returns the marshalled data []byte and an error if any issues occur during the marshaling process.
func MarshalData(data interface{}) ([]byte, error) {
	marshalledData, err := json.Marshal(data)
	if err != nil {
		return nil, goErrorHandler.OperationFailure("marshal data", err)
	}
	return marshalledData, nil
}

// CompareHashAndPassword compares a hashed password with a plaintext password.
// It takes a hashed password and a plaintext password as input and returns an error.
// If the passwords match, the error is nil; otherwise, an error is returned.
func CompareHashAndPassword(hashedPassword string, plaintextPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	if err != nil {
		return goErrorHandler.OperationFailure("compare hash and password", err)
	}
	return nil
}

// HashPassword generates a hashed password using bcrypt with a default cost.
// It takes a plaintext password as input and returns the hashed password as a string.
// An error is returned if the hashing operation fails.
func HashPassword(plaintextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", goErrorHandler.OperationFailure("hash password", err)
	}
	return string(hashedPassword), nil
}

type services struct {
	Redis  *redis.Client       // Redis client
	Rabbit *amqp091.Connection // RabbitMQ connection
	Mongo  *mongo.Client       // MongoDB client
}

// InitializeServices initializes various services such as Redis, RabbitMQ, and MongoDB.
func InitializeServices(ctx context.Context, cfg *config.Config) services {
	// Connect to RabbitMQ server using the provided URL.
	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.AppCfg.RabbitUrl)
	if err != nil {
		log.Println(err) // Log error if RabbitMQ connection fails
		os.Exit(1)       // Exit program if RabbitMQ connection fails
	}
	// Establish a Redis connection
	redisConn := connectToRedis(cfg.AppCfg.RedisUrl)

	// Connect to Mongo database
	mongoClient := connectToMongo(cfg.MongoCfg, ctx)

	return services{
		Redis:  redisConn,   // Assign Redis connection to services
		Rabbit: rabbitConn,  // Assign RabbitMQ connection to services
		Mongo:  mongoClient, // Assign MongoDB client to services
	}
}

// connectToRedis establishes a connection to a Redis server and returns a Redis client.
// It takes the address of the Redis server as a parameter.
func connectToRedis(addr string) *redis.Client {
	// Create a new Redis client with specified options
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		WriteTimeout: 5 * time.Second, // Maximum time to wait for write operations
		ReadTimeout:  5 * time.Second, // Maximum time to wait for read operations
		DialTimeout:  3 * time.Second, // Maximum time to wait for a connection to be established
		MaxRetries:   3,               // Maximum number of retries before giving up on a command
	})
}

// connectToMongo establishes a connection to MongoDB using the provided configuration and context.
func connectToMongo(mongoCfg config.MongoCfg, ctx context.Context) *mongo.Client {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoCfg.MongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: mongoCfg.MongoUsername,     // Set username for authentication
		Password: mongoCfg.MongoUserPassword, // Set password for authentication
	})

	// connect
	c, err := mongo.Connect(ctx, clientOptions) // Establish connection to MongoDB
	if err != nil {
		log.Printf("failed connect to mongo - %v\n", err) // Log error if connection fails
		os.Exit(1)                                        // Exit program if connection fails
	}

	log.Println("Connected to mongo!") // Log successful connection

	return c // Return MongoDB client
}
