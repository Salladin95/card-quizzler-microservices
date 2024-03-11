package server

import (
	cacheManager "github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/handlers"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/middlewares"
	"github.com/Salladin95/rmqtools"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// setupRoutes configures and defines API routes for the Echo server.
func (app *App) setupRoutes(
	broker rmqtools.MessageBroker,
	handlers handlers.ApiHandlersInterface,
	cacheManager cacheManager.CacheManager,
) {
	appThrottler := middlewares.NewThrottler(10, 14)

	routes := app.server.Group("/v1/api")
	routes.Use(appThrottler.Middleware)

	// ********** PROMETHEUS METRICS ****************
	routes.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// ****************** AUTH **********************
	authRoutes := routes.Group("/auth")
	authThrottler := middlewares.NewThrottler(1, 1)
	authRoutes.Use(authThrottler.Middleware)
	authRoutes.POST("/sign-in", handlers.SignIn)
	authRoutes.POST("/sign-up", handlers.SignUp)
	authRoutes.GET("/refresh", handlers.Refresh)

	// ****************** EMAIL *********************
	routes.POST("/request-email-verification", handlers.RequestEmailVerification)
	// ****************** RESET PASSWORD *********************
	routes.PATCH("/user/reset-password", handlers.ResetPassword)
	// ***************** PROTECTED ROUTES ************
	protectedRoutes := routes.Group("")
	protectedRoutes.Use(
		middlewares.AccessTokenValidator(broker, cacheManager, app.config.JwtCfg.JWTAccessSecret),
	)
	// ****************** PROTECTED >> USER *********************
	protectedRoutes.GET("/user/profile", handlers.GetProfile)
	protectedRoutes.GET("/user/:id", handlers.GetUserById)
	protectedRoutes.PATCH("/user/update-email/:id", handlers.UpdateEmail)
	protectedRoutes.PATCH("/user/update-password/:id", handlers.UpdatePassword)
	// ****************** PROTECTED >> CARD-QUIZZLER *********************
	// updates user's progress
	protectedRoutes.PATCH("/process-quiz-result", handlers.ProcessQuizResult)
	// adds folder to user
	protectedRoutes.PATCH("/add-folder-to-user", handlers.AddFolderToUser)
	// creates folder
	protectedRoutes.POST("/folder", handlers.CreateFolder)
	// updates folder
	protectedRoutes.PATCH("/folder/:id", handlers.UpdateFolder)
	// gets folders by userID
	protectedRoutes.GET("/folder", handlers.GetUserFolders)
	// gets folder by folderID
	protectedRoutes.GET("/folder/:id", handlers.GetFolderByID)
	// deletes folder by folderID
	protectedRoutes.DELETE("/folder/:id", handlers.DeleteFolder)
	// adds module to user
	protectedRoutes.PATCH("/add-module-to-user", handlers.AddModuleToUser)
	// adds module to user
	protectedRoutes.PATCH("/add-module-to-folder", handlers.AddModuleToFolder)
	// creates module and adds to the folder with passed folderID
	protectedRoutes.POST("/module/:id", handlers.CreateModuleInFolder)
	// creates a module
	protectedRoutes.POST("/module", handlers.CreateModule)
	// updates module
	protectedRoutes.PATCH("/module/:id", handlers.UpdateModule)
	// get modules by userID
	protectedRoutes.GET("/module", handlers.GetUserModules)
	// get difficult modules by userID
	protectedRoutes.GET("/difficult-modules", handlers.GetDifficultModules)
	// gets module by moduleID
	protectedRoutes.GET("/module/:id", handlers.GetModuleByID)
	// deletes module by moduleID
	protectedRoutes.DELETE("/module/:id", handlers.DeleteModule)
	// deletes module from the folder
	protectedRoutes.PATCH("/delete-module-from-folder", handlers.DeleteModuleFromFolder)
}
