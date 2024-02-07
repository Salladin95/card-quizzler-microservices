package server

import (
	cacheManager "github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/handlers"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/middlewares"
)

// setupRoutes configures and defines API routes for the Echo server.
func (app *App) setupRoutes() {
	cacheManager := cacheManager.NewCacheManager(app.redis, app.config)
	brokerHandlers := handlers.NewHandlers(app.config, app.rabbit, cacheManager)
	routes := app.server.Group("/v1/api")

	// ****************** AUTH **********************
	routes.POST("/auth/sign-in", brokerHandlers.SignIn)
	routes.POST("/auth/sign-up", brokerHandlers.SignUp)
	// ****************** REFRESH *********************
	refreshRoute := routes.Group(
		"/user-service/refresh",
		middlewares.RefreshTokenValidator(cacheManager, app.config.JwtCfg.JWTRefreshSecret),
	)
	refreshRoute.GET("", brokerHandlers.Refresh)
	// ***************** PROTECTED ROUTES ************
	protectedRoutes := routes.Group("")
	protectedRoutes.Use(middlewares.AccessTokenValidator(cacheManager, app.config.JwtCfg.JWTAccessSecret))
	// ****************** PROFILE *********************
	protectedRoutes.GET("/user/profile", brokerHandlers.GetProfile)
	protectedRoutes.GET("/user/:id", brokerHandlers.GetUserById)
	protectedRoutes.GET("/user/", brokerHandlers.GetUsers)
}
