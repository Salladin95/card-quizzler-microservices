package lib

import (
	"encoding/json"
	"github.com/Salladin95/goErrorHandler"
)

// UnmarshalData unmarshals JSON data into the provided unmarshalTo interface.
// It returns an error if any issues occur during the unmarshaling process.
// Note - unmarshalTo must be pointer !!!
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
