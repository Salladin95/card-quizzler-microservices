package entities

import (
	"github.com/Salladin95/goErrorHandler"
	"github.com/go-playground/validator/v10"
)

type SendEmailVerificationDto struct {
	Email string `json:"email" validate:"required"`
}

// Verify validates the structure and content of the SignUpDto.
func (dto *SendEmailVerificationDto) Verify() error {
	// Create a new validator instance.
	validate := validator.New()

	// Validate the SignUpDto structure.
	if err := validate.Struct(dto); err != nil {
		// Convert validation errors and return a ValidationFailure error.
		return goErrorHandler.ValidationFailure(goErrorHandler.ConvertValidationErrors(err))
	}

	return nil
}

type EmailCode struct {
	Email string `json:"email" validate:"required"`
	Code  int    `json:"code" validate:"required"`
}

// Verify validates the structure and content of the SignUpDto.
func (emailCode *EmailCode) Verify() error {
	// Create a new validator instance.
	validate := validator.New()

	// Validate the SignUpDto structure.
	if err := validate.Struct(emailCode); err != nil {
		// Convert validation errors and return a ValidationFailure error.
		return goErrorHandler.ValidationFailure(goErrorHandler.ConvertValidationErrors(err))
	}

	return nil
}

type LogMessage struct {
	FromService string `json:"fromService" validate:"required"`
	Message     string `json:"message" validate:"required"`
	Level       string `json:"level" validate:"required"`
	Name        string `json:"name" validate:"omitempty"`
	Method      string `json:"method" validate:"omitempty"`
}

func (log *LogMessage) GenerateLog(message string, level string, method string, name string) LogMessage {
	return LogMessage{
		Level:       level,
		Method:      method,
		FromService: "mail-service",
		Message:     message,
		Name:        name,
	}
}
