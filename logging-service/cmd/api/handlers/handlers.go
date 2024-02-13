package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/models"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type loggingHandlers struct {
	mongo *mongo.Client
}

type LongingHandlers interface {
	Log(routingKey string, payload []byte)
}

func NewLoggingHandlers(mongo *mongo.Client) LongingHandlers {
	return &loggingHandlers{mongo: mongo}
}

func (lh *loggingHandlers) Log(routingKey string, payload []byte) {
	fmt.Println("Start processing log")
	var logMessage entities.LogMessage
	err := json.Unmarshal(payload, &logMessage)
	if err != nil {
		fmt.Println("failed to unmarshall logMessage")
		return
	}

	err = logMessage.Verify()
	if err != nil {
		fmt.Printf("log-message has failed validation - %v", err)
		return
	}

	fmt.Printf(
		"[logging service][routingkey - %s][message from - %s][method - %s][level - %s][description - %s] message - %s\n",
		routingKey, logMessage.FromService, logMessage.Method, logMessage.Level, logMessage.Name, logMessage.Message,
	)
	// Create a background context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	models := models.NewModels(lh.mongo)

	err = models.LogEntry.Insert(ctx, logMessage)
	if err != nil {
		fmt.Printf("failed to insert log - %v", err)
	}
}
