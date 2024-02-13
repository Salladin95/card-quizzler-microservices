package lib

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/config"
	"github.com/Salladin95/rmqtools"
	"github.com/go-playground/validator/v10"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strings"
)

// ConvertValidationErrors converts validation errors to a more readable format.
func ConvertValidationErrors(err error) map[string]string {
	// Assert that the error is of type validator.ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Log a fatal error if the type assertion fails
		log.Fatal("Unexpected error type during validation")
	}

	// Convert validation errors to a map for easier handling
	validationErrorMap := make(map[string]string)
	for _, fieldError := range validationErrors {
		// Convert field names to lowercase for consistency
		fieldName := strings.ToLower(fieldError.Field())
		// Build a validation error message using the field tag
		validationErrorMap[fieldName] = fmt.Sprintf("Validation failed - %s", fieldError.Tag())
	}

	return validationErrorMap
}

func ValidationFailure(messages map[string]string) error {
	var errorMsg string

	for key, value := range messages {
		errorMsg += fmt.Sprintf("%s: %s\n", key, value)
	}
	return fmt.Errorf(errorMsg)
}

type services struct {
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

	// Connect to Mongo database
	mongoClient := connectToMongo(cfg.MongoCfg, ctx)

	return services{
		Rabbit: rabbitConn,  // Assign RabbitMQ connection to services
		Mongo:  mongoClient, // Assign MongoDB client to services
	}
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
