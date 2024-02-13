package handlers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/constants"
	appEntities "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/cachedRepository"
	userService "github.com/Salladin95/card-quizzler-microservices/user-service/proto"
	"github.com/Salladin95/rmqtools"
)

type UserServer struct {
	userService.UnimplementedUserServiceServer
	Repo   cachedRepository.CachedRepository
	Broker rmqtools.MessageBroker
}

// log sends a log message to the message broker.
func (us *UserServer) log(ctx context.Context, message, level, method string) {
	var log appEntities.LogMessage // Create a new LogMessage struct
	// Push log message to the message broker
	us.Broker.PushToQueue(
		ctx,
		constants.LogCommand, // Specify the log command constant
		// Generate log message with provided details
		log.GenerateLog(message, level, method, "user server, handling gRPC requests"),
	)
}
