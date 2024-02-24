package subscribers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/repositories"
	"github.com/Salladin95/rmqtools"
)

type subscribers struct {
	broker rmqtools.MessageBroker
	repo   repositories.Repository
}

type Subscribers interface {
	SubscribeToUserCreation(ctx context.Context)
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

func (s *subscribers) SubscribeToUserCreation(ctx context.Context) {
	s.log(
		ctx,
		"subscribing to events",
		"info",
		"SubscribeToUserCreation",
	)

	s.broker.ListenForUpdates(
		[]string{constants.CreatedUserKey},
		func(_ string, payload []byte) {
			// TODO: extract to a function
			var createUserDto entities.CreateUserDto
			if err := lib.UnmarshalData(payload, &createUserDto); err != nil {
				s.log(
					ctx,
					fmt.Sprintf("unmarshall payload - %v", err),
					"error",
					"SubscribeToUserCreation",
				)
				return
			}
			if err := createUserDto.Verify(); err != nil {
				s.log(
					ctx,
					fmt.Sprintf("invalid payload - %v", err),
					"error",
					"SubscribeToUserCreation",
				)
				return
			}
			if err := s.repo.CreateUser(createUserDto.ID); err != nil {
				s.log(
					ctx,
					fmt.Sprintf("failed to create user record - %v", err),
					"error",
					"SubscribeToUserCreation",
				)
				return
			}
			s.log(
				ctx,
				"user record is created",
				"info",
				"SubscribeToUserCreation",
			)
		},
	)
}

// log sends a log message to the message broker.
func (s *subscribers) log(ctx context.Context, message, level, method string) {
	var log entities.LogMessage // Create a new LogMessage struct
	// Push log message to the message broker
	s.broker.PushToQueue(
		ctx,
		constants.LogCommand, // Specify the log command constant
		// Generate log message with provided details
		log.GenerateLog(message, level, method, "broker message subscribers"),
	)
}
