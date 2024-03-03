package constants

const (
	AmqpExchange                    = "api-service"
	AmqpQueue                       = "api-service-queue"
	FetchedUserKey                  = "user-service.user.fetched"
	CreatedUserKey                  = "user-service.user.created"
	UpdatedUserKey                  = "user-service.user.updated"
	DeletedUserKey                  = "user-service.user.deleted"
	LogCommand                      = "logging-service.log"
	RequestEmailVerificationCommand = "mail-service.request-email-verification"
	FetchUserFoldersKey             = "card-quizzler-service.user-folders.fetched"
	FetchedFolderKey                = "card-quizzler-service.folder.fetched"
	FetchedUserModulesKey           = "card-quizzler-service.user-modules.fetched"
	FetchedModuleKey                = "card-quizzler-service.module.fetched"
	MutatedFolderKey                = "card-quizzler-service.folder.mutated"
	CreatedFolderKey                = "card-quizzler-service.folder.created"
	CreatedModuleKey                = "card-quizzler-service.module.created"
	MutatedModuleKey                = "card-quizzler-service.module.mutated"
	DeletedFolderKey                = "card-quizzler-service.folder.deleted"
	DeletedModuleKey                = "card-quizzler-service.module.deleted"
)
