package lib

import (
	"encoding/json"
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

// UnmarshalData unmarshals JSON data into the provided unmarshalTo interface.
// It returns an error if any issues occur during the unmarshaling process.
// Not - unmarshalTo must be pointer !!!
func UnmarshalData(data []byte, unmarshalTo interface{}) error {
	err := json.Unmarshal(data, unmarshalTo)
	if err != nil {
		return goErrorHandler.OperationFailure("unmarshal data", err)
	}
	return nil
}

// MarshalData marshals data into a JSON-encoded byte slice.
// It returns the marshalled data []byte and an error if any issues occur during the marshaling process.
func MarshalData(data interface{}) ([]byte, error) {
	marshalledData, err := json.Marshal(data)
	if err != nil {
		return nil, goErrorHandler.OperationFailure("marshal data", err)
	}
	return marshalledData, nil
}