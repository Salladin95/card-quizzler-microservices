package middlewares

import (
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
)

func errorHandler(err error, c echo.Context) {
	apiError := goErrorHandler.MapServiceErrorToAPIError(err)
	c.String(apiError.Status, apiError.Message)
}

func HttpErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			errorHandler(err, c)
		}
		return nil
	}
}
