package lib

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// ConnectToMongo establishes a connection to MongoDB using the provided configuration and context.
func ConnectToMongo(mongoCfg config.MongoCfg, ctx context.Context) *mongo.Client {
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
