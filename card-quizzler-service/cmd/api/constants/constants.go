package constants

const (
	AmqpExchange          = "api-service"
	AmqpQueue             = "api-service-queue"
	LogCommand            = "logging-service.log"
	CreatedUserKey        = "user-service.user.created"
	FetchUserFoldersKey   = "card-quizzler-service.user-folders.fetched"
	FetchFolderKey        = "card-quizzler-service.folder.fetched"
	FetchUserModulesKey   = "card-quizzler-service.user-modules.fetched"
	FetchModuleKey        = "card-quizzler-service.module.fetched"
	MutateFolderKey       = "card-quizzler-service.folder.mutated"
	CreateFolderKey       = "card-quizzler-service.folder.created"
	CreateModuleKey       = "card-quizzler-service.module.created"
	MutateModuleKey       = "card-quizzler-service.module.mutated"
	MutateFolderAndModule = "card-quizzler-service.module&folder.mutated"
	DeleteFolderKey       = "card-quizzler-service.folder.delete"
	DeleteModuleKey       = "card-quizzler-service.module.delete"
	FoldersCacheKey       = "folders"
	FolderCacheKey        = "folder"
	ModulesCacheKey       = "modules"
	ModuleCacheKey        = "module"
)
