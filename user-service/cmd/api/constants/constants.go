package constants

const (
	AmqpExchange = "api-service"
	AmqpQueue    = "api-service-queue"

	SignInKey                    = "user-service.sign-in.command"
	SignUpKey                    = "user-service.sign-up.command"
	FetchedUsersKey              = "user-service.users.fetched"
	FetchedUserKey               = "user-service.user.fetched"
	CreatedUserKey               = "user-service.user.created"
	UpdatedUserKey               = "user-service.user.updated"
	DeletedUserKey               = "user-service.user.deleted"
	LogCommand                   = "logging-service.log"
	EmailVerificationCodeCommand = "mail-service.email-verification-code"
)
