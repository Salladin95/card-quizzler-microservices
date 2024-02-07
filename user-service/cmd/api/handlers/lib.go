package handlers

import (
	"encoding/json"
	user "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/model"
	auth "github.com/Salladin95/card-quizzler-microservices/user-service/proto"
	"github.com/Salladin95/goErrorHandler"
	"net/http"
)

// getErrorStatus maps a Go service error to an API error status code.
// It takes an error as input and returns the corresponding API error status code as an int64.
func getErrorStatus(err error) int64 {
	return int64(goErrorHandler.MapServiceErrorToAPIError(err).Status)
}

// getErrorMessage maps a Go service error to an API error message.
// It takes an error as input and returns the corresponding API error message as a string.
func getErrorMessage(err error) string {
	return goErrorHandler.MapServiceErrorToAPIError(err).Message
}

// buildUserResponse marshals a user and returns an authentication response.
// It takes a user.User pointer, a success code as an int64, and a message as a string.
// The function returns an user-service.Response pointer and an error.
// If marshaling the user fails, it returns a response with a 500 status code and an error message.
func buildUserResponse(u *user.User, successCode int64, message string) (*auth.Response, error) {
	// Marshal the user into JSON format
	marshalledUser, err := json.Marshal(u)
	if err != nil {
		// Return a response with a 500 status code and an error message if marshaling fails
		goErrorHandler.OperationFailure("marshal user", err)
		return &auth.Response{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}
	// Return the authentication response with the success code, message, and marshalled user data
	return &auth.Response{Code: successCode, Message: message, Data: marshalledUser}, nil
}
