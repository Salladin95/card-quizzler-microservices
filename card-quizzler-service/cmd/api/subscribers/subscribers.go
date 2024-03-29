package subscribers

import (
	"context"
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
