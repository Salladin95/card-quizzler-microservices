package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/handlers"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// App represents the main application structure.
type App struct {
	server *echo.Echo       // Echo HTTP server instance
	rabbit *amqp.Connection // RabbitMQ connection instance
	config *config.Config   // Application configuration
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
	lib.Logger.Println("********* START SERVER ***************")
	// Setup middlewares for the Echo server.
	app.setupMiddlewares()

	// Start the Echo server in a goroutine.
	go func() {
		serverAddr := fmt.Sprintf(":%s", app.config.AppCfg.ApiServicePort)
		if err := app.server.Start(serverAddr); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Println("********* SHUTTING DOWN THE SERVER **********")
			lib.Logger.Println("********* SHUTTING DOWN THE SERVER **********")
		}
	}()

	cacheManager := cacheManager.NewCacheManager(app.redis, app.config, app.rabbit)
	handlers := handlers.NewHandlers(app.config, app.rabbit, cacheManager)

	// Setup routes for the Echo server.
	app.setupRoutes(handlers, cacheManager)

	go cacheManager.ListenForUpdates()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Create a context with a timeout for the server shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shut down the Echo server.
	if err := app.server.Shutdown(ctx); err != nil {
		app.server.Logger.Fatal(err)
	}
}
