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

type foldersDto struct {
	UserID string            `json:"userID"`
	Key    string            `json:"key"`
	Data   []entities.Folder `json:"data"`
}

type modulesDto struct {
	UserID string            `json:"userID"`
	Key    string            `json:"key"`
	Data   []entities.Module `json:"data"`
}

// cardQuizEventHandler is a callback function to handle user-related events received from RabbitMQ.
// It processes each event based on the Key and performs corresponding actions.
func (s *subscribers) cardQuizEventHandler(key string, payload []byte) {
	ctx := context.Background()

	s.log(ctx, fmt.Sprintf("start processing Key - %s", key), "info", "cardQuizEventHandler")

	var user userReference
	if err := lib.UnmarshalData(payload, &user); err != nil {
		s.log(ctx, "failed to unmarshall payload to extract userID", "error", "cardQuizEventHandler")
		return
	}

	switch key {
	case constants.FetchUserFoldersKey:
		var dto foldersDto
		if err := lib.UnmarshalData(payload, &dto); err != nil {
			return
		}
		marshalledData, err := lib.MarshalData(dto.Data)
		if err != nil {
			return
		}

		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			cacheManager.FoldersKey(user.UserID),
			dto.Key,
			marshalledData,
			s.cacheManager.Exp(),
		)
		s.log(
			ctx,
			"fetched folders case",
			"info",
			"cardQuizEventHandler",
		)

	case constants.FetchedUserModulesKey:
		var dto modulesDto
		if err := lib.UnmarshalData(payload, &dto); err != nil {
			return
		}
		marshalledData, err := lib.MarshalData(dto.Data)
		if err != nil {
			return
		}
		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			cacheManager.ModulesKey(user.UserID),
			dto.Key,
			marshalledData,
			s.cacheManager.Exp(),
		)
		s.log(
			ctx,
			"fetched modules case",
			"info",
			"cardQuizEventHandler",
		)
	case constants.FetchedDifficultModulesKey:
		var modules []entities.Module
		if err := lib.UnmarshalData(payload, &modules); err != nil {
			return
		}
		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			cacheManager.ModulesKey(user.UserID),
			cacheManager.DifficultModules,
			payload,
			s.cacheManager.Exp(),
		)
		s.log(
			ctx,
			"fetched difficult modules case",
			"info",
			"cardQuizEventHandler",
		)
	case constants.FetchedFolderKey:
		var folder entities.Folder
		if err := lib.UnmarshalData(payload, &folder); err != nil {
			return
		}
		// Clear the cache for the folder
		s.cacheManager.SetCacheByKeys(
			cacheManager.FolderKey(user.UserID),
			cacheManager.FolderKey(folder.ID.String()),
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
		var module entities.Module
		if err := lib.UnmarshalData(payload, &module); err != nil {
			return
		}
		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			cacheManager.ModuleKey(user.UserID),
			cacheManager.ModuleKey(module.ID.String()),
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
		s.cacheManager.ClearCacheByKey(
			cacheManager.FoldersKey(user.UserID),
		)
		s.log(
			ctx,
			"new folder case, clearing cache for folders",
			"info",
			"cardQuizEventHandler",
		)

	case constants.CreatedModuleKey:
		// Clear the cache for the modules
		s.cacheManager.ClearCacheByKey(
			cacheManager.ModulesKey(user.UserID),
		)
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
		if err := s.cacheManager.ClearCacheByKey(
			cacheManager.FoldersKey(user.UserID),
		); err != nil {
			return
		}
		if err := s.cacheManager.ClearCacheByKeys(
			cacheManager.FolderKey(user.UserID),
			cacheManager.FolderKey(folder.ID.String()),
		); err != nil {
			return
		}
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

		// if module has folders, clean there cache as well
		var foldersIDS []string
		for _, folder := range module.Folders {
			foldersIDS = append(foldersIDS, folder.ID.String())
		}

		if len(foldersIDS) > 0 {
			if err := s.cacheManager.ClearCacheByKey(
				cacheManager.FoldersKey(user.UserID),
			); err != nil {
				return
			}
			for _, folderID := range foldersIDS {
				if err := s.cacheManager.ClearCacheByKeys(
					cacheManager.FolderKey(user.UserID),
					cacheManager.FolderKey(folderID),
				); err != nil {
					return
				}
			}

		}
		if err := s.cacheManager.ClearCacheByKey(
			cacheManager.ModulesKey(user.UserID),
		); err != nil {
			return
		}

		if err := s.cacheManager.ClearCacheByKeys(
			cacheManager.ModuleKey(user.UserID),
			cacheManager.ModuleKey(module.ID.String()),
		); err != nil {
			return
		}

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
