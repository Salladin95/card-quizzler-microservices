package lib

import (
	"github.com/Salladin95/goErrorHandler"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// DataWithVerify is an interface that requires a Verify method.
type DataWithVerify interface {
	Verify() error
}

// Verify validates the given structure
func Verify(data interface{}) error {
	// Create a new validator instance.
	validate := validator.New()

	// Validate the SignUpDto structure.
	if err := validate.Struct(data); err != nil {
		// Convert validation errors and return a ValidationFailure error.
		return goErrorHandler.ValidationFailure(goErrorHandler.ConvertValidationErrors(err))
	}

	return nil
}

// BindBodyAndVerify binds the request body to a DataWithVerify interface
// and then calls the Verify method on the provided data.
// !note - data must be a pointer
func BindBodyAndVerify(c echo.Context, data DataWithVerify) error {
	// Bind the request body to the DataWithVerify interface
	if err := c.Bind(data); err != nil {
		return goErrorHandler.BindRequestToBodyFailure(err)
	}

	// Call the Verify method on the provided data
	err := data.Verify()
	return err
}
