package subscribers

import (
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
	lib.LogInfo(fmt.Sprintf("[cardQuizEventHandler] Start processing Key - %s", key))

	var user userReference
	if err := lib.UnmarshalData(payload, &user); err != nil {
		lib.LogError(fmt.Errorf("[cardQuizEventHandler] Failed to unmarshall payload to extract userID"))
		return
	}

	switch key {
	case constants.FetchUserFoldersKey:
		var dto foldersDto
		if err := lib.UnmarshalData(payload, &dto); err != nil {
			lib.LogError(err)
			return
		}
		marshalledData, err := lib.MarshalData(dto.Data)
		if err != nil {
			lib.LogError(err)
			return
		}

		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			cacheManager.FoldersKey(user.UserID),
			dto.Key,
			marshalledData,
			s.cacheManager.Exp(),
		)

	case constants.FetchedUserModulesKey:
		var dto modulesDto
		if err := lib.UnmarshalData(payload, &dto); err != nil {
			lib.LogError(err)
			return
		}
		marshalledData, err := lib.MarshalData(dto.Data)
		if err != nil {
			lib.LogError(err)
			return
		}
		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			cacheManager.ModulesKey(user.UserID),
			dto.Key,
			marshalledData,
			s.cacheManager.Exp(),
		)
	case constants.FetchedDifficultModulesKey:
		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			cacheManager.ModulesKey(user.UserID),
			cacheManager.DifficultModules,
			payload,
			s.cacheManager.Exp(),
		)
	case constants.FetchedFolderKey:
		var folder entities.Folder
		if err := lib.UnmarshalData(payload, &folder); err != nil {
			lib.LogError(err)
			return
		}
		// Clear the cache for the folder
		s.cacheManager.SetCacheByKeys(
			cacheManager.FolderKey(user.UserID),
			cacheManager.FolderKey(folder.ID.String()),
			payload,
			s.cacheManager.Exp(),
		)

	case constants.FetchedModuleKey:
		var module entities.Module
		if err := lib.UnmarshalData(payload, &module); err != nil {
			lib.LogError(err)
			return
		}
		// Clear the cache for the folders
		s.cacheManager.SetCacheByKeys(
			cacheManager.ModuleKey(user.UserID),
			cacheManager.ModuleKey(module.ID.String()),
			payload,
			s.cacheManager.Exp(),
		)
	case constants.CreatedFolderKey:
		// Clear the cache for the folders
		s.cacheManager.ClearCacheByKey(
			cacheManager.FoldersKey(user.UserID),
		)

	case constants.CreatedModuleKey:
		// Clear the cache for the modules
		s.cacheManager.ClearCacheByKey(
			cacheManager.ModulesKey(user.UserID),
		)

	case constants.MutatedFolderKey, constants.DeletedFolderKey:
		// Clear the cache for the folders
		var folder entities.Folder
		if err := lib.UnmarshalData(payload, &folder); err != nil {
			lib.LogError(err)
			return
		}
		s.cacheManager.ClearCacheByKey(
			cacheManager.FoldersKey(user.UserID),
		)
		s.cacheManager.ClearCacheByKeys(
			cacheManager.FolderKey(user.UserID),
			cacheManager.FolderKey(folder.ID.String()),
		)

	case constants.MutatedModuleKey, constants.DeletedModuleKey:
		// Clear the cache for the folders
		var module entities.Module
		if err := lib.UnmarshalData(payload, &module); err != nil {
			lib.LogError(err)
			return
		}

		s.cacheManager.ClearCacheByKey(
			cacheManager.ModulesKey(user.UserID),
		)

		s.cacheManager.ClearCacheByKeys(
			cacheManager.ModuleKey(user.UserID),
			cacheManager.ModuleKey(module.ID.String()),
		)

		// if module has folders, clean there cache as well
		var foldersIDS []string
		for _, folder := range module.Folders {
			foldersIDS = append(foldersIDS, folder.ID.String())
		}

		if len(foldersIDS) > 0 {
			s.cacheManager.ClearCacheByKey(
				cacheManager.FoldersKey(user.UserID),
			)
			for _, folderID := range foldersIDS {
				s.cacheManager.ClearCacheByKeys(
					cacheManager.FolderKey(user.UserID),
					cacheManager.FolderKey(folderID),
				)
			}

		}
	default:
		lib.LogInfo("unknown case")
	}
}
