package subscribers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/shared"
)

func (s *subscribers) subscribeToUserCreation(ctx context.Context) {
	s.broker.ListenForUpdates(
		[]string{constants.CreatedUserKey},
		func(_ string, payload []byte) {
			var createUserDto entities.CreateUserDto
			if err := lib.UnmarshalData(payload, &createUserDto); err != nil {
				lib.LogError(
					fmt.Errorf("unmarshall payload - %v", err),
				)
				return
			}
			if err := createUserDto.Verify(); err != nil {
				lib.LogError(
					fmt.Errorf("invalid payload - %v", err),
				)
				return
			}
			if err := s.repo.CreateUser(createUserDto.ID); err != nil {
				lib.LogError(
					fmt.Errorf("failed to create user record - %v", err),
				)
				return
			}
		},
	)
}
