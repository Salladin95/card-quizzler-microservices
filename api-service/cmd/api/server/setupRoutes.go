package server

import (
	cacheManager "github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/handlers"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/messageBroker"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/middlewares"
)

// setupRoutes configures and defines API routes for the Echo server.
func (app *App) setupRoutes(
	broker messageBroker.MessageBroker,
	handlers handlers.BrokerHandlersInterface,
	cacheManager cacheManager.CacheManager,
) {
	routes := app.server.Group("/v1/api")
	// ****************** AUTH **********************
	routes.POST("/auth/sign-in", handlers.SignIn)
	routes.POST("/auth/sign-up", handlers.SignUp)
	// ****************** REFRESH *********************
	refreshRoute := routes.Group(
		"/user-service/refresh",
		middlewares.RefreshTokenValidator(broker, cacheManager, app.config.JwtCfg.JWTRefreshSecret),
	)
	refreshRoute.GET("", handlers.Refresh)
	// ***************** PROTECTED ROUTES ************
	protectedRoutes := routes.Group("")
	protectedRoutes.Use(middlewares.AccessTokenValidator(broker, cacheManager, app.config.JwtCfg.JWTAccessSecret))
	// ****************** PROFILE *********************
	protectedRoutes.GET("/user/profile", handlers.GetProfile)
	protectedRoutes.GET("/user/:id", handlers.GetUserById)
	protectedRoutes.GET("/user/", handlers.GetUsers)
}
