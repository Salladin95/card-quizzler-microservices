package messageBroker

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
)

func (mb *messageBroker) GenerateLogEvent(ctx context.Context, logMessage entities.LogMessage) error {
	return mb.PushToQueue(ctx, constants.LogCommand, logMessage)
}
