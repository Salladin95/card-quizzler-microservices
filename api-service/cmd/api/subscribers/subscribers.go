package subscribers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/rmqtools"
	"log"
)

type subscribers struct {
	broker       rmqtools.MessageBroker
	cacheManager cacheManager.CacheManager
}

type Subscribers interface {
	SubscribeToUserEvents(ctx context.Context)
}

func NewMessageBrokerSubscribers(
	broker rmqtools.MessageBroker,
	cacheManager cacheManager.CacheManager,
) Subscribers {
	return &subscribers{
		broker:       broker,
		cacheManager: cacheManager,
	}
}

var userEvents = []string{
	constants.CreatedUserKey,
	constants.UpdatedUserKey,
	constants.DeletedUserKey,
	constants.FetchedUserKey,
}

func (s *subscribers) SubscribeToUserEvents(ctx context.Context) {
	s.log(
		ctx,
		"subscribing to user events",
		"info",
		"SubscribeToUserEvents",
	)

	e := s.broker.ListenForUpdates(userEvents, s.userEventHandler)
	s.log(ctx, e.Error(), "error", "SubscribeToUserEvents")
}

// userEventHandler is a callback function to handle user-related events received from RabbitMQ.
// It processes each event based on the key and performs corresponding actions.
func (s *subscribers) userEventHandler(key string, payload []byte) {
	ctx := context.Background()

	s.log(ctx, fmt.Sprintf("start processing key - %s", key), "info", "ListenForUpdates")

	var user entities.UserResponse
	err := lib.UnmarshalData(payload, &user)
	if err != nil {
		msg := fmt.Sprintf("user event handler failed to unmarshall user - %v", err)
		s.log(ctx, msg, "error", "ListenForUpdates")
		log.Panic(msg)
		return
	}

	switch key {
	case constants.CreatedUserKey:
		s.cacheManager.SetCacheInPipeline(
			s.cacheManager.UserHashKey(user.ID),
			cacheManager.UserKey,
			payload,
			s.cacheManager.Exp(),
		)
		// Clear the cache for the user list
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UsersKey)

		s.log(
			ctx,
			"new user case, clearing cache for [s.key, email, id]",
			"error",
			"ListenForUpdates",
		)
	case constants.UpdatedUserKey:
		// Clear the cache
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UsersKey)
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UserKey)

		s.log(
			ctx,
			"user updated case, clearing cache for [s.key, email, id]",
			"error",
			"ListenForUpdates",
		)
	case constants.DeletedUserKey:
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UsersKey)
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UserKey)

		s.log(
			ctx,
			"user updated case, clearing cache for [s.key, email, id]",
			"error",
			"ListenForUpdates",
		)
	case constants.FetchedUserKey:
		s.cacheManager.SetCacheInPipeline(
			s.cacheManager.UserHashKey(user.ID),
			cacheManager.UserKey,
			payload,
			s.cacheManager.Exp(),
		)

		s.log(
			ctx,
			"fetched user case, setting cache",
			"info",
			"ListenForUpdates",
		)
	default:
		s.log(
			ctx,
			"unknown case",
			"error",
			"ListenForUpdates",
		)
	}
}

func (s *subscribers) log(ctx context.Context, message, level, method string) {
	var log entities.LogMessage

	// Push log message to the message broker
	s.broker.PushToQueue(
		ctx,
		constants.LogCommand, // Specify the log command constant
		// Generate log message with provided details
		log.GenerateLog(message, level, method, "user events listener"),
	)
}
