package handlers

import "github.com/Salladin95/goErrorHandler"

func getErrorStatus(err error) int64 {
	return int64(goErrorHandler.MapServiceErrorToAPIError(err).Status)
}

func getErrorMessage(err error) string {
	return goErrorHandler.MapServiceErrorToAPIError(err).Message
}
