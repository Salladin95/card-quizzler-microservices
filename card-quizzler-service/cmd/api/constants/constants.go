package constants

const (
	AmqpExchange          = "api-service"
	AmqpQueue             = "api-service-queue"
	LogCommand            = "logging-service.log"
	CreatedUserKey        = "user-service.user.created"
	FetchedUserFoldersKey = "card-quizzler-service.user-folders.fetched"
	FetchedFolderKey      = "card-quizzler-service.folder.fetched"
	FetchedUserModulesKey = "card-quizzler-service.user-modules.fetched"
	FetchedModuleKey      = "card-quizzler-service.module.fetched"
	MutatedFolderKey      = "card-quizzler-service.folder.mutated"
	CreatedFolderKey      = "card-quizzler-service.folder.created"
	CreatedModuleKey      = "card-quizzler-service.module.created"
	MutatedModuleKey      = "card-quizzler-service.module.mutated"
	DeletedFolderKey      = "card-quizzler-service.folder.deleted"
	DeletedModuleKey      = "card-quizzler-service.module.deleted"
	FoldersCacheKey       = "folders"
	FolderCacheKey        = "folder"
	ModulesCacheKey       = "modules"
	ModuleCacheKey        = "module"
)
