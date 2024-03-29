package subscribers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/cachedRepository"
	"github.com/Salladin95/rmqtools"
)

type subscribers struct {
	broker     rmqtools.MessageBroker
	cachedRepo cachedRepository.CachedRepository
}

type Subscribers interface {
	SubscribeToEmailVerificationReqs(ctx context.Context)
}

func NewMessageBrokerSubscribers(
	broker rmqtools.MessageBroker,
	cachedRepo cachedRepository.CachedRepository,
) Subscribers {
	return &subscribers{
		broker:     broker,
		cachedRepo: cachedRepo,
	}
}

func (s *subscribers) SubscribeToEmailVerificationReqs(ctx context.Context) {
	s.broker.ListenForUpdates(
		[]string{constants.EmailVerificationCodeCommand},
		func(_ string, payload []byte) {
			if err := s.cachedRepo.SetEmailVerificationCode(ctx, payload); err != nil {
				lib.LogError(err.Error())
			}
		},
	)
}
