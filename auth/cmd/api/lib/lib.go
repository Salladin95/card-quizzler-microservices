package lib

import (
	"encoding/json"
	"github.com/Salladin95/goErrorHandler"
	"golang.org/x/crypto/bcrypt"
)

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

// CompareHashAndPassword compares a hashed password with a plaintext password.
// It takes a hashed password and a plaintext password as input and returns an error.
// If the passwords match, the error is nil; otherwise, an error is returned.
func CompareHashAndPassword(hashedPassword string, plaintextPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	if err != nil {
		return goErrorHandler.OperationFailure("compare has and password", err)
	}
	return nil
}

// HashPassword generates a hashed password using bcrypt with a default cost.
// It takes a plaintext password as input and returns the hashed password as a string.
// An error is returned if the hashing operation fails.
func HashPassword(plaintextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", goErrorHandler.OperationFailure("hash password", err)
	}
	return string(hashedPassword), nil
}
