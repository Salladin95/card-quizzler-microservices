package lib

import (
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
)

// BindBody bind body from request to bindTo param
// Note - bindTo must be pointer !!!
func BindBody(c echo.Context, bindTo interface{}) error {
	// Bind the request body to the DataWithVerify interface
	if err := c.Bind(&bindTo); err != nil {
		return goErrorHandler.BindRequestToBodyFailure(err)
	}
	return nil
}

// BindBodyAndVerify binds the request body to a DataWithVerify interface
// and then calls the Verify method on the provided data.
// Note - bindTo must be pointer !!!
func BindBodyAndVerify(c echo.Context, bindTo DataWithVerify) error {
	// Bind the request body to the DataWithVerify interface
	if err := c.Bind(bindTo); err != nil {
		return goErrorHandler.BindRequestToBodyFailure(err)
	}

	// Call the Verify method on the provided bindTo
	err := bindTo.Verify()
	return err
}
