package handlers

import (
	"encoding/json"
	user "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/model"
	auth "github.com/Salladin95/card-quizzler-microservices/auth-service/proto"
	"github.com/Salladin95/goErrorHandler"
	"net/http"
)

func getErrorStatus(err error) int64 {
	return int64(goErrorHandler.MapServiceErrorToAPIError(err).Status)
}

func getErrorMessage(err error) string {
	return goErrorHandler.MapServiceErrorToAPIError(err).Message
}

// Marshals user and returns auth response
func buildUserResponse(u *user.User, successCode int64, message string) (*auth.Response, error) {
	marshalledUer, err := json.Marshal(u)
	if err != nil {
		goErrorHandler.OperationFailure("marshal user", err)
		return &auth.Response{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}
	return &auth.Response{Code: successCode, Message: message, Data: marshalledUer}, nil
}
