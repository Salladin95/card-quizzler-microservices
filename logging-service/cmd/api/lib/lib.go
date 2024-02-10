package lib

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"strings"
)

// ConvertValidationErrors converts validation errors to a more readable format.
func ConvertValidationErrors(err error) map[string]string {
	// Assert that the error is of type validator.ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Log a fatal error if the type assertion fails
		log.Fatal("Unexpected error type during validation")
	}

	// Convert validation errors to a map for easier handling
	validationErrorMap := make(map[string]string)
	for _, fieldError := range validationErrors {
		// Convert field names to lowercase for consistency
		fieldName := strings.ToLower(fieldError.Field())
		// Build a validation error message using the field tag
		validationErrorMap[fieldName] = fmt.Sprintf("Validation failed - %s", fieldError.Tag())
	}

	return validationErrorMap
}

func ValidationFailure(messages map[string]string) error {
	var errorMsg string

	for key, value := range messages {
		errorMsg += fmt.Sprintf("%s: %s\n", key, value)
	}
	return fmt.Errorf(errorMsg)
}
