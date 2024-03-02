package subscribers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
)

func (s *subscribers) subscribeToUserCreation(ctx context.Context) {
	s.log(
		ctx,
		"subscribing to events",
		"info",
		"Listen",
	)

	s.broker.ListenForUpdates(
		[]string{constants.CreatedUserKey},
		func(_ string, payload []byte) {
			var createUserDto entities.CreateUserDto
			if err := lib.UnmarshalData(payload, &createUserDto); err != nil {
				s.log(
					ctx,
					fmt.Sprintf("unmarshall payload - %v", err),
					"error",
					"subscribeToUserCreation",
				)
				return
			}
			if err := createUserDto.Verify(); err != nil {
				s.log(
					ctx,
					fmt.Sprintf("invalid payload - %v", err),
					"error",
					"subscribeToUserCreation",
				)
				return
			}
			if err := s.repo.CreateUser(createUserDto.ID); err != nil {
				s.log(
					ctx,
					fmt.Sprintf("failed to create user record - %v", err),
					"error",
					"subscribeToUserCreation",
				)
				return
			}
			s.log(
				ctx,
				"user record is created",
				"info",
				"subscribeToUserCreation",
			)
		},
	)
}
