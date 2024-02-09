package handlers

import (
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/cachedRepository"
	userService "github.com/Salladin95/card-quizzler-microservices/user-service/proto"
)

type UserServer struct {
	userService.UnimplementedUserServiceServer
	Repo cachedRepository.CachedRepository
}
