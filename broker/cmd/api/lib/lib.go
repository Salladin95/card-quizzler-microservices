package lib

import (
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
)

// DataWithVerify is an interface that requires a Verify method.
type DataWithVerify interface {
	Verify() error
}

// BindBodyToStruct binds the request body to the given data structure.
func BindBodyToStruct(c echo.Context, data any) error {
	// Bind the request body to the provided data structure
	if err := c.Bind(data); err != nil {
		return goErrorHandler.BindRequestToBodyFailure(err)
	}
	return nil
}

// BindBodyAndVerify binds the request body to a DataWithVerify interface
// and then calls the Verify method on the provided data.
func BindBodyAndVerify(c echo.Context, data DataWithVerify) error {
	// Bind the request body to the DataWithVerify interface
	if err := c.Bind(&data); err != nil {
		return goErrorHandler.BindRequestToBodyFailure(err)
	}

	// Call the Verify method on the provided data
	err := data.Verify()
	return err
}
