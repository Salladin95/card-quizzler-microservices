package cacheManager

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/rmqtools"
)

// ListenForUpdates listens for updates from RabbitMQ and handles user-related events.
// It creates a new consumer for the specified AMQP exchange and queue,
// then listens for events with the provided keys and handles them using the handleUserEvents method.
func (cm *cacheManager) ListenForUpdates() {
	consumer, err := rmqtools.NewConsumer(
		cm.rabbitConn,
		constants.AmqpExchange,
		constants.AmqpQueue,
	)
	if err != nil {
		lib.Logger.Error().Err(err)
	}
	err = consumer.Listen(
		[]string{
			constants.CreatedUserKey,
			constants.UpdatedUserKey,
			constants.DeletedUserKey,
			constants.FetchedUserKey,
			constants.FetchedUsersKey,
		},
		cm.handleUserEvents,
	)
	if err != nil {
		lib.Logger.Error().Err(err)
	}
}

// handleUserEvents is a callback function to handle user-related events received from RabbitMQ.
// It processes each event based on the key and performs corresponding actions.
func (cm *cacheManager) handleUserEvents(key string, payload []byte) {
	ctx := context.Background()
	cm.pushToQueue(ctx, constants.LogCommand, generateLog(fmt.Sprintf("start processing key - %s", key), "info"))
	switch key {
	case constants.CreatedUserKey:
		cm.pushToQueue(ctx, constants.LogCommand, generateLog("new user case, clearing cache for [cm.userKey, email, id]", "info"))
		handleUserCache(cm, payload)
		// Clear the cache for the user list
		cm.ClearDataByKey(cm.userKey)
	case constants.UpdatedUserKey:
		cm.pushToQueue(ctx, constants.LogCommand, generateLog("user updated case, clearing cache for [cm.userKey, email, id]", "info"))
		handleUserCache(cm, payload)
		// Clear the cache for the user list
		cm.ClearDataByKey(cm.userKey)
	case constants.DeletedUserKey:
		cm.pushToQueue(ctx, constants.LogCommand, generateLog("user deleted case, clearing cache for [cm.userKey, email, id]", "info"))
		user, err := entities.UnmarshalUser(payload)
		if err != nil {
			cm.redisClient.FlushAll()
			cm.pushToQueue(ctx, constants.LogCommand, generateLog("failed to unmarshal user, flushing cache", "error"))
			return
		}
		cm.ClearDataByKey(cm.userHashKey(user.ID.String()))
		cm.ClearDataByKey(user.Email)
	case constants.FetchedUserKey:
		cm.pushToQueue(ctx, constants.LogCommand, generateLog("fetched users case, setting cache", "info"))
		handleUserCache(cm, payload)
	case constants.FetchedUsersKey:
		cm.pushToQueue(ctx, constants.LogCommand, generateLog("fetched user case, setting cache", "info"))
		cm.SetCacheByKey(cm.userKey, payload)
	default:
		cm.pushToQueue(ctx, constants.LogCommand, generateLog("unknown case", "error"))
	}
}

// handleUserCache parses user data from the payload, caches it, and handles any errors.
// It unmarshals the user data, caches it using both the hash key derived from the user ID
// and the email as cache keys, and flushes the cache if unmarshaling fails.
func handleUserCache(cm *cacheManager, payload []byte) {
	user, err := entities.UnmarshalUser(payload)

	if err != nil {
		lib.Logger.Info().Msg("**** failed to unmarshal user, flushing cache ********")
		lib.Logger.Error().Err(err)
	}

	// Cache the newly created user using both the hash key derived from the user ID and the email as cache keys
	cm.SetCacheByKey(cm.userHashKey(user.ID.String()), payload)
	cm.SetCacheByKey(user.Email, payload)
}

// TODO: ADD CONSISTENT WAY OF LOGGING
func generateLog(message string, level string) entities.LogMessage {
	return entities.LogMessage{
		Level:       level,
		Method:      "handleUserEvents",
		FromService: "api-service",
		Message:     message,
		Name:        "handler events for rabbitMQ consumer",
	}
}
