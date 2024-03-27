package subscribers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/repositories"
	"github.com/Salladin95/rmqtools"
)

type subscribers struct {
	broker rmqtools.MessageBroker
	repo   repositories.Repository
}

type Subscribers interface {
	Listen(ctx context.Context)
}

func (s *subscribers) Listen(ctx context.Context) {
	s.subscribeToUserCreation(ctx)
}

func NewMessageBrokerSubscribers(
	broker rmqtools.MessageBroker,
	repo repositories.Repository,
) Subscribers {
	return &subscribers{
		broker: broker,
		repo:   repo,
	}
}

// log sends a log message to the message broker.
func (s *subscribers) log(ctx context.Context, message, level, method string) {
	var log entities.LogMessage // Create a new LogMessage struct
	// Push log message to the message broker
	if err := s.broker.PushToQueue(
		ctx,
		constants.LogCommand, // Specify the log command constant
		// Generate log message with provided details
		log.GenerateLog(message, level, method, "broker message subscribers"),
	); err != nil {
		fmt.Printf("[subscribers] Failed to push log - %v\n", err)
	}
}
