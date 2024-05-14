package entities

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/Salladin95/goErrorHandler"
)

type JsonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CreateUserDto struct {
	ID string `json:"id" validate:"required"`
}

type ResultTerm struct {
	models.Term
	ID     string `json:"id" validate:"required"`
	Answer bool   `json:"answer" validate:"required"`
}

type SecureAccess struct {
	Access   models.AccessType `json:"access" validate:"required"`
	Password string            `json:"password" validate:"omitempty,min=4"`
}

func ValidateSecureAccess(dto *SecureAccess) error {
	switch dto.Access {
	case models.AccessOpen, models.AccessOnlyMe:
		return nil
	case models.AccessPassword:
		if dto.Password == "" {
			return goErrorHandler.NewError(goErrorHandler.ErrBadRequest, fmt.Errorf("password is required"))
		}
		return nil
	default:
		return goErrorHandler.NewError(goErrorHandler.ErrBadRequest, fmt.Errorf("invalid access type"))
	}
}

func (dto *SecureAccess) Verify() error {
	if err := lib.Verify(dto); err != nil {
		return lib.Verify(dto)
	}
	return ValidateSecureAccess(dto)
}
