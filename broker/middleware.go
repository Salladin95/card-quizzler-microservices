package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *App) useMiddleware() {
	// specify who is allowed to connect
	app.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO: FIX IN PRODUCTION - ["https://*", "http://*"]
		AllowOrigins: []string{"https://*", "http://*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXCSRFToken, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		MaxAge:       300,
	}))
}
