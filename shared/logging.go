package lib

import (
	"github.com/labstack/echo/v4"
	"log/slog"
)

func LogInfo(msg string, args ...any) {
	slog.Info(msg, args...)
}

func LogError(err error, args ...any) {
	slog.Error(err.Error(), args...)
}

func LogRequestInfo(c echo.Context, msg string) {
	slog.Info(
		msg,
		slog.String("path", c.Request().URL.Path),
		slog.String("method", c.Request().Method),
	)
}

func LogRequestError(c echo.Context, msg string) {
	slog.Error(
		msg,
		slog.String("path", c.Request().URL.Path),
		slog.String("method", c.Request().Method),
	)
}
