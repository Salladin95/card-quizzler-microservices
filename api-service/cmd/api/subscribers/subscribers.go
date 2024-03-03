package subscribers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/rmqtools"
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
	constants.FetchedFolderKey,
	constants.FetchedModuleKey,
	constants.DeletedFolderKey,
	constants.DeletedModuleKey,
	constants.MutatedFolderKey,
	constants.MutatedModuleKey,
}

func (s *subscribers) Listen(ctx context.Context) {
	s.log(
		ctx,
		"subscribing events",
		"info",
		"Listen",
	)

	go s.listenToUserEvents(ctx)
	go s.listenToCardQuizEvents(ctx)
}

func (s *subscribers) listenToUserEvents(ctx context.Context) {
	s.log(
		ctx,
		"subscribing to user events",
		"info",
		"listenToUserEvents",
	)

	e := s.broker.ListenForUpdates(userEvents, s.userEventHandler)
	s.log(ctx, e.Error(), "error", "listenToUserEvents")
}

func (s *subscribers) listenToCardQuizEvents(ctx context.Context) {
	s.log(
		ctx,
		"subscribing to card quiz events",
		"info",
		"listenToCardQuizEvents",
	)

	e := s.broker.ListenForUpdates(cardQuizEvents, s.cardQuizEventHandler)
	s.log(ctx, e.Error(), "error", "listenToCardQuizEvents")
}

func (s *subscribers) log(ctx context.Context, message, level, method string) {
	var log entities.LogMessage

	// Push log message to the message broker
	s.broker.PushToQueue(
		ctx,
		constants.LogCommand, // Specify the log command constant
		// Generate log message with provided details
		log.GenerateLog(message, level, method, "events listener"),
	)
}
