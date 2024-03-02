package subscribers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"log"
)

// userEventHandler is a callback function to handle user-related events received from RabbitMQ.
// It processes each event based on the key and performs corresponding actions.
func (s *subscribers) userEventHandler(key string, payload []byte) {
	ctx := context.Background()

	s.log(ctx, fmt.Sprintf("start processing key - %s", key), "info", "userEventHandler")

	var user entities.UserResponse
	err := lib.UnmarshalData(payload, &user)
	if err != nil {
		msg := fmt.Sprintf("user event handler failed to unmarshall user - %v", err)
		s.log(ctx, msg, "error", "userEventHandler")
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
			"userEventHandler",
		)
	case constants.UpdatedUserKey:
		// Clear the cache
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UsersKey)
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UserKey)

		s.log(
			ctx,
			"user updated case, clearing cache for [s.key, email, id]",
			"error",
			"userEventHandler",
		)
	case constants.DeletedUserKey:
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UsersKey)
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UserKey)

		s.log(
			ctx,
			"user updated case, clearing cache for [s.key, email, id]",
			"error",
			"userEventHandler",
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
			"userEventHandler",
		)
	default:
		s.log(
			ctx,
			"unknown case",
			"error",
			"userEventHandler",
		)
	}
}
