package server

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/middlewares"
	"github.com/Salladin95/rmqtools"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// setupMiddlewares configures and adds middlewares to the Echo server.
func (app *App) setupMiddlewares(broker rmqtools.MessageBroker) {
	// Use the custom HTTP error handler middleware.
	app.server.Use(middlewares.HttpErrorHandler(broker))

	// Configure CORS middleware to control cross-origin resource sharing.
	app.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO: FIX IN PRODUCTION - AllowOrigins should be limited to specific domains.
		// Specify the allowed origins. In production, replace ["https://*", "http://*"] with specific domain(s).
		AllowOrigins: []string{"https://*", "http://*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderXCSRFToken,
			echo.HeaderAuthorization,
		},
		MaxAge: 300, // Set the maximum age for preflight requests.
	}))
}
