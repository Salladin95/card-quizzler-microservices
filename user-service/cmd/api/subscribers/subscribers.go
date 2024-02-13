package subscribers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/constants"
	appEntities "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/entities"
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
	s.log(
		ctx,
		"subscribing to user verification requests",
		"info",
		"SubscribeToEmailVerificationReqs",
	)

	s.broker.ListenForUpdates(
		[]string{constants.EmailVerificationCodeCommand},
		func(_ string, payload []byte) {
			err := s.cachedRepo.SetEmailVerificationCode(ctx, payload)
			if err != nil {
				s.log(ctx, err.Error(), "error", "SubscribeToEmailVerificationReqs")
			}
			s.log(
				ctx,
				"email verification code is saved to cache",
				"info",
				"SubscribeToEmailVerificationReqs",
			)
		},
	)
}

// log sends a log message to the message broker.
func (s *subscribers) log(ctx context.Context, message, level, method string) {
	var log appEntities.LogMessage // Create a new LogMessage struct
	// Push log message to the message broker
	s.broker.PushToQueue(
		ctx,
		constants.LogCommand, // Specify the log command constant
		// Generate log message with provided details
		log.GenerateLog(message, level, method, "broker message subscribers"),
	)
}
