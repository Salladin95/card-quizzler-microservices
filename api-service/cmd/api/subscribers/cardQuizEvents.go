package subscribers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
)

type userReference struct {
	UserID string `json:"userID"`
}

// cardQuizEventHandler is a callback function to handle user-related events received from RabbitMQ.
// It processes each event based on the key and performs corresponding actions.
func (s *subscribers) cardQuizEventHandler(key string, payload []byte) {
	ctx := context.Background()

	s.log(ctx, fmt.Sprintf("start processing key - %s", key), "info", "cardQuizEventHandler")

	var user userReference
	if key != constants.FetchUserFoldersKey && key != constants.FetchedUserModulesKey {
		if err := lib.UnmarshalData(payload, &user); err != nil {
			s.log(ctx, "failed to unmarshall payload to extract userID", "error", "cardQuizEventHandler")
			return
		}
	}

	switch key {
	case constants.FetchUserFoldersKey:
		var folders []entities.Folder
		if err := lib.UnmarshalData(payload, &folders); err != nil || len(folders) < 1 {
			return
		}
		uid := folders[0].UserID
		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			s.cacheManager.UserHashKey(uid),
			cacheManager.Folders,
			payload,
			s.cacheManager.Exp(),
		)
		s.log(
			ctx,
			"fetched folders case",
			"info",
			"cardQuizEventHandler",
		)

	case constants.FetchedUserModulesKey:
		var modules []entities.Module
		if err := lib.UnmarshalData(payload, &modules); err != nil || len(modules) < 1 {
			return
		}
		uid := modules[0].UserID
		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			s.cacheManager.UserHashKey(uid),
			cacheManager.Modules,
			payload,
			s.cacheManager.Exp(),
		)
		s.log(
			ctx,
			"fetched modules case",
			"info",
			"cardQuizEventHandler",
		)
	case constants.FetchedFolderKey:
		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			s.cacheManager.UserHashKey(user.UserID),
			cacheManager.Folder,
			payload,
			s.cacheManager.Exp(),
		)
		s.log(
			ctx,
			"fetched folder case",
			"info",
			"cardQuizEventHandler",
		)

	case constants.FetchedModuleKey:
		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			s.cacheManager.UserHashKey(user.UserID),
			cacheManager.Module,
			payload,
			s.cacheManager.Exp(),
		)
		s.log(
			ctx,
			"fetched module case",
			"info",
			"cardQuizEventHandler",
		)
	case constants.CreatedFolderKey:
		// Clear the cache for the folders
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.UserID), cacheManager.Folders)
		s.log(
			ctx,
			"new folder case, clearing cache for folders",
			"info",
			"cardQuizEventHandler",
		)

	case constants.CreatedModuleKey:
		// Clear the cache for the folders
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(user.UserID), cacheManager.Modules)
		s.log(
			ctx,
			"new module case, clearing cache for modules",
			"info",
			"cardQuizEventHandler",
		)

	case constants.MutatedFolderKey, constants.DeletedFolderKey:
		// Clear the cache for the folders
		var folder entities.Folder
		if err := lib.UnmarshalData(payload, &folder); err != nil {
			return
		}
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(folder.UserID), cacheManager.Folders)
		s.cacheManager.ClearCacheByKeys(
			s.cacheManager.UserHashKey(folder.UserID),
			fmt.Sprintf(
				"%s-%s",
				cacheManager.Folder,
				folder.ID.String(),
			),
		)
		s.log(
			ctx,
			"[mutated, deleted] folder case, clearing cache for folders",
			"info",
			"cardQuizEventHandler",
		)

	case constants.MutatedModuleKey, constants.DeletedModuleKey:
		// Clear the cache for the folders
		var module entities.Module
		if err := lib.UnmarshalData(payload, &module); err != nil {
			return
		}
		s.cacheManager.ClearCacheByKeys(s.cacheManager.UserHashKey(module.UserID), cacheManager.Modules)
		s.cacheManager.ClearCacheByKeys(
			s.cacheManager.UserHashKey(module.UserID),
			fmt.Sprintf(
				"%s-%s",
				cacheManager.Module,
				module.ID.String(),
			),
		)
		s.log(
			ctx,
			"[mutated, deleted] module case, clearing cache for modules",
			"info",
			"cardQuizEventHandler",
		)
	default:
		s.log(
			ctx,
			"unknown case",
			"error",
			"cardQuizEventHandler",
		)
	}
}
