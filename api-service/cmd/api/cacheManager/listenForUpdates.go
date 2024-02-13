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
	cm.broker.ListenForUpdates(userEvents, cm.userEventHandler)
}

// userEventHandler is a callback function to handle user-related events received from RabbitMQ.
// It processes each event based on the key and performs corresponding actions.
func (cm *cacheManager) userEventHandler(key string, payload []byte) {
	ctx := context.Background()

	cm.log(ctx, fmt.Sprintf("start processing key - %s", key), "info", "ListenForUpdates")

	// payload should type of entities.UserResponse in other cases
	// so we unmarshal it and use it id to interact with cache
	var user entities.UserResponse
	if key != constants.FetchedUsersKey {
		err := lib.UnmarshalData(payload, &user)

		if err != nil {
			msg := fmt.Sprintf("user event handler failed to unmarshall user - %v", err)
			cm.log(ctx, msg, "error", "ListenForUpdates")
			log.Panic(msg)
		}
	}

	switch key {
	case constants.CreatedUserKey:
		cm.setCacheInPipeline(cm.userHashKey(user.ID), userKey, payload, cm.exp)
		// Clear the cache for the user list
		cm.ClearCacheByKeys(cm.userHashKey(user.ID), usersKey)

		cm.log(
			ctx,
			"new user case, clearing cache for [cm.key, email, id]",
			"error",
			"ListenForUpdates",
		)
	case constants.UpdatedUserKey:
		cm.setCacheInPipeline(cm.userHashKey(user.ID), userKey, payload, cm.exp)
		// Clear the cache for the user list
		cm.ClearCacheByKeys(cm.userHashKey(user.ID), usersKey)

		cm.log(
			ctx,
			"user updated case, clearing cache for [cm.key, email, id]",
			"error",
			"ListenForUpdates",
		)
	case constants.DeletedUserKey:
		cm.ClearCacheByKeys(cm.userHashKey(user.ID), usersKey)
		cm.ClearCacheByKeys(cm.userHashKey(user.ID), userKey)

		cm.log(
			ctx,
			"user updated case, clearing cache for [cm.key, email, id]",
			"error",
			"ListenForUpdates",
		)
	case constants.FetchedUserKey:
		cm.setCacheInPipeline(cm.userHashKey(user.ID), userKey, payload, cm.exp)

		cm.log(
			ctx,
			"fetched user case, setting cache",
			"info",
			"ListenForUpdates",
		)
	case constants.FetchedUsersKey:
		cm.setCacheInPipeline(cm.userHashKey(user.ID), usersKey, payload, cm.exp)

		cm.log(
			ctx,
			"fetched users case, setting cache",
			"info",
			"ListenForUpdates",
		)
	default:
		cm.log(
			ctx,
			"unknown case",
			"error",
			"ListenForUpdates",
		)
	}
}
