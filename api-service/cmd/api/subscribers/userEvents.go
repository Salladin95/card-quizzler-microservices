package subscribers

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"log"
)

// userEventHandler is a callback function to handle user-related events received from RabbitMQ.
// It processes each event based on the Key and performs corresponding actions.
func (s *subscribers) userEventHandler(key string, payload []byte) {
	lib.LogInfo(fmt.Sprintf("[userEventHandler] Start processing Key - %s", key))

	var user entities.UserResponse
	err := lib.UnmarshalData(payload, &user)
	if err != nil {
		msg := fmt.Sprintf("user event handler failed to unmarshall user - %v", err)
		lib.LogError(err)
		log.Panic(msg)
		return
	}

	switch key {
	case constants.CreatedUserKey:
		s.cacheManager.SetCacheByKeys(
			s.cacheManager.UserHashKey(user.ID),
			cacheManager.UserKey,
			payload,
			s.cacheManager.Exp(),
		)
		// Clear the cache for the user list
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UsersKey)

	case constants.UpdatedUserKey:
		// Clear the cache
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UsersKey)
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UserKey)

	case constants.DeletedUserKey:
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UsersKey)
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.ID), cacheManager.UserKey)

	case constants.FetchedUserKey:
		s.cacheManager.SetCacheByKeys(
			s.cacheManager.UserHashKey(user.ID),
			cacheManager.UserKey,
			payload,
			s.cacheManager.Exp(),
		)

	default:
		lib.LogInfo("unknown case")
	}
}
