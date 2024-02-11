package middlewares

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/messageBroker"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
)

// errorHandler maps a service-specific error to an API error and sends the appropriate response.
func errorHandler(err error, c echo.Context, broker messageBroker.MessageBroker) {
	// Map the service-specific error to an API error.
	apiError := goErrorHandler.MapServiceErrorToAPIError(err)

	var logMessage entities.LogMessage
	broker.GenerateLogEvent(
		c.Request().Context(),
		logMessage.GenerateLog(
			apiError.Message,
			"error",
			"errorHandler",
			"http error handler",
		),
	)

	// Send the API error response with the corresponding HTTP status code and message.
	c.String(apiError.Status, apiError.Message)
}

// HttpErrorHandler is a middleware that catches errors from subsequent middleware or handlers
// and uses the errorHandler function to send an appropriate API error response.
func HttpErrorHandler(broker messageBroker.MessageBroker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Call the next middleware or handler and catch any errors that occur.
			if err := next(c); err != nil {
				// If an error occurs, handle it using the errorHandler function.
				errorHandler(err, c, broker)
			}
			return nil
		}
	}
}
