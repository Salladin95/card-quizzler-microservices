package entities

import (
	"github.com/Salladin95/goErrorHandler"
	"github.com/go-playground/validator/v10"
)

// SignInDto represents the data transfer object for user sign-in requests.
type SignInDto struct {
	Email    string `json:"email" validate:"required,email"`    // Email field with validation rules
	Password string `json:"password" validate:"min=6,required"` // Password field with validation rules
}

// Verify validates the structure and content of the SignInDto.
func (signInDto *SignInDto) Verify() error {
	// Create a new validator instance.
	validate := validator.New()

	// Validate the SignInDto structure.
	if err := validate.Struct(signInDto); err != nil {
		// Convert validation errors and return a ValidationFailure error.
		return goErrorHandler.ValidationFailure(goErrorHandler.ConvertValidationErrors(err))
	}

	return nil
}

// SignUpDto represents the data transfer object for user sign-up requests.
type SignUpDto struct {
	Name     string `json:"name" validate:"required,min=1"`     // Name field with validation rules
	Password string `json:"password" validate:"required,min=6"` // Password field with validation rules
	Email    string `json:"email" validate:"required,email"`    // Email field with validation rules
	Birthday string `json:"birthday" validate:"required,min=1"` // Birthday field with validation rules
}

// Verify validates the structure and content of the SignUpDto.
func (signUpDto *SignUpDto) Verify() error {
	// Create a new validator instance.
	validate := validator.New()

	// Validate the SignUpDto structure.
	if err := validate.Struct(signUpDto); err != nil {
		// Convert validation errors and return a ValidationFailure error.
		return goErrorHandler.ValidationFailure(goErrorHandler.ConvertValidationErrors(err))
	}

	return nil
}

// JsonResponse represents a simple JSON response message structure.
type JsonResponse struct {
	Message string `json:"message"` // Message field for JSON response messages
}
