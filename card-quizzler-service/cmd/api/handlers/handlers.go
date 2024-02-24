package handlers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/repositories"
	quizService "github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/proto"
	"github.com/Salladin95/rmqtools"
	"google.golang.org/grpc"
)

type CardQuizzlerServer struct {
	quizService.UnimplementedCardQuizzlerServiceServer
	repo   repositories.Repository
	broker rmqtools.MessageBroker
}

func RegisterQuizzlerServer(
	gRPCServer *grpc.Server,
	repo repositories.Repository,
	broker rmqtools.MessageBroker,
) {
	// Register the AuthServer implementation with the gRPC server.
	quizService.RegisterCardQuizzlerServiceServer(
		gRPCServer,
		&CardQuizzlerServer{
			repo:   repo,
			broker: broker,
		},
	)
}

// log sends a log message to the message broker.
func (cq *CardQuizzlerServer) log(ctx context.Context, message, level, method string) {
	var log entities.LogMessage // Create a new LogMessage struct
	// Push log message to the message broker
	cq.broker.PushToQueue(
		ctx,
		constants.LogCommand, // Specify the log command constant
		// Generate log message with provided details
		log.GenerateLog(message, level, method, "card-quiz-server, handling gRPC requests"),
	)
}
