package server

import (
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// setupMiddlewares configures and adds middlewares to the Echo server.
func (app *App) setupMiddlewares() {
	// Use the custom HTTP error handler middleware.
	app.server.Use(middlewares.HttpErrorHandler)

	// Configure CORS middleware to control cross-origin resource sharing.
	app.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO: FIX IN PRODUCTION - AllowOrigins should be limited to specific domains.
		// Specify the allowed origins. In production, replace ["https://*", "http://*"] with specific domain(s).
		AllowOrigins: []string{"https://*", "http://*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXCSRFToken, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		MaxAge:       300, // Set the maximum age for preflight requests.
	}))
}
