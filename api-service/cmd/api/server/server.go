package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/handlers"
	broker "github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/messageBroker"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// App represents the main application structure.
type App struct {
	config *config.Config   // Application configuration
	server *echo.Echo       // Echo HTTP server instance
	rabbit *amqp.Connection // RabbitMQ connection instance
	redis  *redis.Client    // Redis client
}

// IApp defines the interface for the main application.
type IApp interface {
	Start()
}

// NewApp creates a new instance of the application.
func NewApp(cfg *config.Config, rabbit *amqp.Connection, redisClient *redis.Client) IApp {
	return &App{
		server: echo.New(),
		rabbit: rabbit,
		config: cfg,
		redis:  redisClient,
	}
}

// Start initializes and starts the application.
func (app *App) Start() {
	mBroker := broker.NewMessageBroker(app.rabbit)
	cacheManager := cacheManager.NewCacheManager(app.redis, app.config, mBroker)
	handlers := handlers.NewHandlers(app.config, mBroker, cacheManager)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	mBroker.GenerateLogEvent(
		ctx,
		generateServerLog(
			fmt.Sprintf("start the server on the port - %s", app.config.AppCfg.ApiServicePort),
			"info",
		),
	)
	// Setup middlewares for the Echo server.
	app.setupMiddlewares(mBroker)
	// Setup routes for the Echo server.
	app.setupRoutes(mBroker, handlers, cacheManager)

	// Start the Echo server in a goroutine.
	go func() {
		serverAddr := fmt.Sprintf(":%s", app.config.AppCfg.ApiServicePort)
		if err := app.server.Start(serverAddr); err != nil && errors.Is(err, http.ErrServerClosed) {
			mBroker.GenerateLogEvent(
				ctx,
				generateServerLog("shutting down the server", "info"),
			)
		}
	}()

	go cacheManager.ListenForUpdates()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Gracefully shut down the Echo server.
	if err := app.server.Shutdown(ctx); err != nil {
		mBroker.GenerateLogEvent(
			ctx,
			generateServerLog(err.Error(), "error"),
		)
		app.server.Logger.Fatal(err)
	}
}

func generateServerLog(message string, level string) entities.LogMessage {
	var logMessage entities.LogMessage
	return logMessage.GenerateLog(message, level, "START", "setting up server")
}
