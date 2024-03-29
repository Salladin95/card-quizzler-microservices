package subscribers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/rmqtools"
	"log/slog"
)

type subscribers struct {
	broker       rmqtools.MessageBroker
	cacheManager cacheManager.CacheManager
}

type Subscribers interface {
	Listen(ctx context.Context)
}

func NewMessageBrokerSubscribers(
	broker rmqtools.MessageBroker,
	cacheManager cacheManager.CacheManager,
) Subscribers {
	return &subscribers{
		broker:       broker,
		cacheManager: cacheManager,
	}
}

var userEvents = []string{
	constants.CreatedUserKey,
	constants.UpdatedUserKey,
	constants.DeletedUserKey,
	constants.FetchedUserKey,
}

var cardQuizEvents = []string{
	constants.CreatedFolderKey,
	constants.CreatedModuleKey,
	constants.FetchUserFoldersKey,
	constants.FetchedUserModulesKey,
	constants.FetchedDifficultModulesKey,
	constants.FetchedFolderKey,
	constants.FetchedModuleKey,
	constants.DeletedFolderKey,
	constants.DeletedModuleKey,
	constants.MutatedFolderKey,
	constants.MutatedModuleKey,
}

func (s *subscribers) Listen(ctx context.Context) {
	go s.listenToUserEvents(ctx)
	go s.listenToCardQuizEvents(ctx)
}

func (s *subscribers) listenToUserEvents(ctx context.Context) {
	lib.LogInfo(
		"subscribing to user events",
	)

	if e := s.broker.ListenForUpdates(userEvents, s.userEventHandler); e != nil {
		lib.LogError(e)
	}
}

func (s *subscribers) listenToCardQuizEvents(ctx context.Context) {
	slog.Info(
		"subscribing to card quiz events",
	)

	if e := s.broker.ListenForUpdates(cardQuizEvents, s.cardQuizEventHandler); e != nil {
		lib.LogError(e)
	}
}
