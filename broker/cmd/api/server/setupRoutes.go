package server

import (
	cacheManager "github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/handlers"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/middlewares"
	"github.com/labstack/echo/v4"
	"net/http"
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
		"/auth/refresh",
		middlewares.RefreshTokenValidator(cacheManager, app.config.JwtCfg.JWTRefreshSecret),
	)
	refreshRoute.GET("", brokerHandlers.Refresh)
	// ***************** PROTECTED ROUTES ************
	protectedRoutes := routes.Group("")
	protectedRoutes.Use(middlewares.AccessTokenValidator(cacheManager, app.config.JwtCfg.JWTAccessSecret))
	// ****************** PROFILE *********************
	// TODO: REPLACE MOCK HANDLER
	protectedRoutes.GET("/user/profile", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface {
		}{
			"message": "here we go again",
		})
	})
}
