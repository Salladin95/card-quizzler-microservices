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
	handlers handlers.ApiHandlersInterface,
	cacheManager cacheManager.CacheManager,
) {
	routes := app.server.Group("/v1/api")

	// ****************** AUTH **********************
	authRoutes := routes.Group("/auth")
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
	protectedRoutes.POST("/process-quiz-result", handlers.ProcessQuizResult)
	// adds folder to user
	protectedRoutes.POST("/add-folder-to-user", handlers.AddFolderToUser)
	// creates folder
	protectedRoutes.POST("/folder", handlers.CreateFolder)
	// updates folder
	protectedRoutes.PATCH("/folder", handlers.UpdateFolder)
	// gets folders by userID
	protectedRoutes.GET("/folders/:uid", handlers.GetUserFolders)
	// gets folder by folderID
	protectedRoutes.GET("/folder/:id", handlers.GetFolderByID)
	// deletes folder by folderID
	protectedRoutes.DELETE("/folder/:id", handlers.DeleteFolder)
	// adds module to user
	protectedRoutes.POST("/add-module-to-user", handlers.AddModuleToUser)
	// adds module to user
	protectedRoutes.POST("/add-module-to-folder", handlers.AddModuleToFolder)
	// creates module and adds to the folder with passed folderID
	protectedRoutes.POST("/module/:id", handlers.CreateModuleInFolder)
	// creates a module
	protectedRoutes.POST("/module", handlers.CreateModule)
	// updates module
	protectedRoutes.PATCH("/module/:id", handlers.UpdateModule)
	// gets module by moduleID
	protectedRoutes.GET("/module/:id", handlers.GetModuleByID)
	// get modules by userID
	protectedRoutes.GET("/modules/:uid", handlers.GetUserModules)
	// deletes module by moduleID
	protectedRoutes.DELETE("/module/:id", handlers.DeleteModule)
	// deletes module from the folder
	protectedRoutes.DELETE("/delete-module-from-folder", handlers.DeleteModuleFromFolder)
}
