package middlewares

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/rmqtools"
)

func logTokenValidation(broker rmqtools.MessageBroker, message, level, method string) {
	ctx := context.Background()
	broker.PushToQueue(
		ctx,
		constants.LogCommand,
		generateTokenValidatorLog(
			message,
			level,
			method,
		),
	)
}

func generateTokenValidatorLog(message string, level string, method string) entities.LogMessage {
	var logMessage entities.LogMessage
	return logMessage.GenerateLog(message, level, method, "access token validtor")
}
