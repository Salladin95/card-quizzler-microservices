package appEntities

import (
	"github.com/Salladin95/goErrorHandler"
	"github.com/go-playground/validator/v10"
)

// JsonResponse represents a simple JSON response message structure.
type JsonResponse struct {
	Message string `json:"message"` // Message field for JSON response messages
}

type EmailCode struct {
	Email string `json:"email" validate:"required"`
	Code  int    `json:"code" validate:"required"`
}

// Verify validates the structure and content of the SignInDto.
func (payload *EmailCode) Verify() error {
	// Create a new validator instance.
	validate := validator.New()

	// Validate the SignInDto structure.
	if err := validate.Struct(payload); err != nil {
		// Convert validation errors and return a ValidationFailure error.
		return goErrorHandler.ValidationFailure(goErrorHandler.ConvertValidationErrors(err))
	}

	return nil
}
