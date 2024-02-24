package handlers

import (
	"encoding/json"
	quizService "github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/proto"
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

// buildSuccessfulResponse marshals a user and returns a quizService. response.
// It takes a user.User pointer, a success code as an int64, and a message as a string.
// The function returns a quizService.Response pointer and an error.
// If marshaling the user fails, it returns a response with a 500 status code and an error message.
func buildSuccessfulResponse(data interface{}, successCode int64, message string) (*quizService.Response, error) {
	// Marshal the user into JSON format
	marshalledData, err := json.Marshal(data)
	if err != nil {
		return &quizService.Response{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}
	// Return the authentication response with the success code, message, and marshalled user data
	return &quizService.Response{Code: successCode, Message: message, Data: marshalledData}, nil
}

// buildFailedResponse extracts code and message from error and returns a quizService.Response.
func buildFailedResponse(err error) (*quizService.Response, error) {
	// Extract status code and message from the error
	return &quizService.Response{Code: getErrorStatus(err), Message: getErrorMessage(err)}, nil
}
