package server

import (
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/handlers"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/middlewares"
	"github.com/labstack/echo/v4"
	"net/http"
)

// setupRoutes configures and defines API routes for the Echo server.
func (app *App) setupRoutes() {
	routes := app.server.Group("/v1/api")
	brokerHandlers := handlers.NewHandlers(app.config, app.rabbit)

	// ****************** AUTH **********************
	routes.POST("/auth/sign-in", brokerHandlers.SignIn)
	routes.POST("/auth/sign-up", brokerHandlers.SignUp)
	// ****************** REFRESH *********************
	refreshRoute := routes.Group(
		"/auth/refresh",
		middlewares.RefreshTokenValidator(app.config.JwtCfg.JWTRefreshSecret),
	)
	refreshRoute.GET("", brokerHandlers.Refresh)
	// ***************** PROTECTED ROUTES ************
	protectedRoutes := routes.Group("")
	protectedRoutes.Use(middlewares.AccessTokenValidator(app.config.JwtCfg.JWTAccessSecret))
	// ****************** PROFILE *********************
	// TODO: REPLACE MOCK HANDLER
	protectedRoutes.GET("/user/profile", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface {
		}{
			"message": "here we go again",
		})
	})
}
