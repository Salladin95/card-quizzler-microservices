package server

import (
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/handlers"
)

// setupRoutes configures and defines API routes for the Echo server.
func (app *App) setupRoutes() {
	// Create a group of routes with the "/v1/api" prefix.
	routes := app.server.Group("/v1/api")
	// Initialize handlers for the API routes.
	bHandlers := handlers.NewHandlers(app.config, app.rabbit)

	// ****************** AUTH **********************
	// Define a route for user sign-in.
	routes.POST("/auth/sign-in", bHandlers.SignIn)
	// Define a route for user sign-up.
	routes.POST("/auth/sign-up", bHandlers.SignUp)
}
