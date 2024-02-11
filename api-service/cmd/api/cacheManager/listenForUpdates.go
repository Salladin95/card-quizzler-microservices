package cacheManager

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"log"
)

var userEvents = []string{
	constants.CreatedUserKey,
	constants.UpdatedUserKey,
	constants.DeletedUserKey,
	constants.FetchedUserKey,
	constants.FetchedUsersKey,
}

// ListenForUpdates listens for updates from RabbitMQ and handles user-related events.
// It creates a new consumer for the specified AMQP exchange and queue,
// then listens for events with the provided keys and handles them using the userEventHandler method.
func (cm *cacheManager) ListenForUpdates() {
	cm.messageBroker.ListenForUpdates(userEvents, cm.userEventHandler)
}

// userEventHandler is a callback function to handle user-related events received from RabbitMQ.
// It processes each event based on the key and performs corresponding actions.
func (cm *cacheManager) userEventHandler(key string, payload []byte) {
	ctx := context.Background()
	cm.messageBroker.GenerateLogEvent(
		ctx,
		generateUserEventHandlerLog(
			fmt.Sprintf("start processing key - %s", key),
			"info",
		),
	)

	// payload should type of entities.UserResponse in other cases
	// so we unmarshal it and use it id to interact with cache
	var user entities.UserResponse
	if key != constants.FetchedUsersKey {
		err := lib.UnmarshalData(payload, &user)

		if err != nil {
			log.Panic(fmt.Errorf("user event handler failed to unmarshall user - %v", err))
		}
	}

	switch key {
	case constants.CreatedUserKey:
		cm.messageBroker.GenerateLogEvent(
			ctx,
			generateUserEventHandlerLog(
				"new user case, clearing cache for [cm.key, email, id]",
				"info",
			),
		)
		cm.setCacheInPipeline(cm.userHashKey(user.ID.String()), userKey, payload, cm.exp)
		// Clear the cache for the user list
		cm.ClearCacheByKeys(cm.userHashKey(user.ID.String()), usersKey)
	case constants.UpdatedUserKey:
		cm.messageBroker.GenerateLogEvent(
			ctx,
			generateUserEventHandlerLog(
				"user updated case, clearing cache for [cm.key, email, id]",
				"info",
			),
		)
		cm.setCacheInPipeline(cm.userHashKey(user.ID.String()), userKey, payload, cm.exp)
		// Clear the cache for the user list
		cm.ClearCacheByKeys(cm.userHashKey(user.ID.String()), usersKey)
	case constants.DeletedUserKey:
		cm.messageBroker.GenerateLogEvent(
			ctx,
			generateUserEventHandlerLog(
				"user deleted case, clearing cache for [cm.key, email, id]",
				"info",
			),
		)
		cm.ClearCacheByKeys(cm.userHashKey(user.ID.String()), usersKey)
		cm.ClearCacheByKeys(cm.userHashKey(user.ID.String()), userKey)
	case constants.FetchedUserKey:
		cm.messageBroker.GenerateLogEvent(
			ctx,
			generateUserEventHandlerLog("fetched user case, setting cache", "info"),
		)
		cm.setCacheInPipeline(cm.userHashKey(user.ID.String()), userKey, payload, cm.exp)
	case constants.FetchedUsersKey:
		cm.messageBroker.GenerateLogEvent(
			ctx,
			generateUserEventHandlerLog("fetched users case, setting cache", "info"),
		)
		cm.setCacheInPipeline(cm.userHashKey(user.ID.String()), usersKey, payload, cm.exp)
	default:
		cm.messageBroker.GenerateLogEvent(
			ctx,
			generateUserEventHandlerLog("unknown case", "error"),
		)
	}
}

// generateUserEventHandlerLog - creates logs for userEventHandler function
func generateUserEventHandlerLog(message string, level string) entities.LogMessage {
	var logMessage *entities.LogMessage
	return logMessage.GenerateLog(
		message,
		level,
		"userEventHandler",
		"handler events for rabbitMQ consumer",
	)
}
