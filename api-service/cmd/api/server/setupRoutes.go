package server

import (
	cacheManager "github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/handlers"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/middlewares"
	"github.com/Salladin95/rmqtools"
)

// setupRoutes configures and defines API routes for the Echo server.
func (app *App) setupRoutes(
	broker rmqtools.MessageBroker,
	handlers handlers.BrokerHandlersInterface,
	cacheManager cacheManager.CacheManager,
) {
	routes := app.server.Group("/v1/api")

	// ****************** AUTH **********************
	authRoutes := routes.Group("/auth")
	authRoutes.POST("/sign-in", handlers.SignIn)
	authRoutes.POST("/sign-up", handlers.SignUp)
	authRoutes.GET("/refresh", handlers.Refresh)

	// ***************** PROTECTED ROUTES ************
	protectedRoutes := routes.Group("")
	protectedRoutes.Use(
		middlewares.AccessTokenValidator(broker, cacheManager, app.config.JwtCfg.JWTAccessSecret),
	)

	// ****************** PROTECTED >> PROFILE *********************
	protectedRoutes.GET("/user/profile", handlers.GetProfile)
	protectedRoutes.GET("/user/:id", handlers.GetUserById)
	protectedRoutes.PATCH("/user/update-email/:id", handlers.UpdateEmail)
	// ****************** PROTECTED >> EMAIL *********************
	protectedRoutes.GET("/request-email-verification", handlers.RequestEmailVerification)
}
