package entities

import (
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/lib"
	"github.com/go-playground/validator/v10"
)

type LogMessage struct {
	FromService string `json:"fromService" validate:"required"`
	Message     string `json:"message" validate:"required"`
	Level       string `json:"level" validate:"required"`
	Name        string `json:"name" validate:"omitempty"`
	Method      string `json:"method" validate:"omitempty"`
}

func (dto *LogMessage) Verify() error {
	// Create a new validator instance.
	validate := validator.New()

	// Validate the SignInDto structure.
	if err := validate.Struct(dto); err != nil {
		// Convert validation errors and return a ValidationFailure error.
		return lib.ValidationFailure(lib.ConvertValidationErrors(err))
	}

	return nil
}
