package handlers

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/repositories"
	quizService "github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/proto"
	"github.com/Salladin95/rmqtools"
	"google.golang.org/grpc"
)

type CardQuizzlerServer struct {
	quizService.UnimplementedCardQuizzlerServiceServer
	Repo   repositories.Repository
	Broker rmqtools.MessageBroker
}

func RegisterQuizzlerServer(
	gRPCServer *grpc.Server,
	repo repositories.Repository,
	broker rmqtools.MessageBroker,
) {
	// Register the CardQuizzlerServer implementation with the gRPC server.
	quizService.RegisterCardQuizzlerServiceServer(
		gRPCServer,
		&CardQuizzlerServer{
			Repo:   repo,
			Broker: broker,
		},
	)
}
