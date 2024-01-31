package main

import (
	"github.com/Salladin95/card-quizzler-microservices/broker-service/handlers"
)

func (app *App) setupRoutes() {
	routes := app.server.Group("/v1/api")
	// initialize handlers
	bHandlers := handlers.NewHandlers(app.rabbit)
	// ****************** AUTH **********************
	routes.POST("/auth/sign-in", bHandlers.SignIn)
	routes.POST("/auth/sign-up", bHandlers.SignUp)
}
