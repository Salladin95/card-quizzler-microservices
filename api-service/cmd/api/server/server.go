package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/handlers"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/subscribers"
	"github.com/Salladin95/rmqtools"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// App represents the main application structure.
type App struct {
	config *config.Config         // Application configuration
	server *echo.Echo             // Echo HTTP server instance
	broker rmqtools.MessageBroker // Message broker interface for communication
	redis  *redis.Client          // Redis client
}

// IApp defines the interface for the main application.
type IApp interface {
	Start()
}

// NewApp creates a new instance of the application.
func NewApp(cfg *config.Config, redisClient *redis.Client, broker rmqtools.MessageBroker) IApp {
	return &App{
		server: echo.New(),
		config: cfg,
		redis:  redisClient,
		broker: broker,
	}
}

// Start initializes and starts the application.
func (app *App) Start() {
	cacheManager := cacheManager.NewCacheManager(app.redis, app.config, app.broker)
	handlers := handlers.NewHandlers(app.config, app.broker, cacheManager)
	listeners := subscribers.NewMessageBrokerSubscribers(app.broker, cacheManager)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	lib.LogInfo(
		fmt.Sprintf("start the server on the port - %s", app.config.AppCfg.ApiServicePort),
	)

	// Setup middlewares for the Echo server.
	app.setupMiddlewares()
	// Setup routes for the Echo server.
	app.setupRoutes(app.broker, handlers, cacheManager)

	// Start the Echo server in a goroutine.
	go func() {
		serverAddr := fmt.Sprintf(":%s", app.config.AppCfg.ApiServicePort)
		if err := app.server.Start(serverAddr); err != nil && errors.Is(err, http.ErrServerClosed) {
			lib.LogInfo("shutting down the server")
		}
	}()

	listeners.Listen(ctx)

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Gracefully shut down the Echo server.
	if err := app.server.Shutdown(ctx); err != nil {
		lib.LogError(err)
		app.server.Logger.Fatal(err)
	}
}
