package user

import (
	"github.com/Salladin95/goErrorHandler"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
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

type UpdateDto struct {
	Name      string     `json:"name" validate:"min=1,omitempty"`
	Email     string     `json:"email" validate:"email,omitempty"`
	Birthday  string     `json:"birthday" validate:"omitempty"`
	Password  string     `json:"password" validate:"omitempty,min=6"`
	CreatedAt time.Time  `json:"createdAt" validate:"omitempty"`
	UpdatedAt time.Time  `json:"updatedAt" validate:"omitempty"`
	DeletedAt *time.Time `json:"deletedAt" validate:"omitempty"`
}

// Verify validates the structure and content of the SignUpDto.
func (updateDto *UpdateDto) Verify() error {
	// Create a new validator instance.
	validate := validator.New()

	// Validate the SignUpDto structure.
	if err := validate.Struct(updateDto); err != nil {
		// Convert validation errors and return a ValidationFailure error.
		return goErrorHandler.ValidationFailure(goErrorHandler.ConvertValidationErrors(err))
	}

	return nil
}

type Response struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Birthday  string     `json:"birthday"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}