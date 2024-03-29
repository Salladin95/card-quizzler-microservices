package handlers

import (
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/cachedRepository"
	userService "github.com/Salladin95/card-quizzler-microservices/user-service/proto"
	"github.com/Salladin95/rmqtools"
)

type UserServer struct {
	userService.UnimplementedUserServiceServer
	CachedRepo cachedRepository.CachedRepository
	Broker     rmqtools.MessageBroker
}
