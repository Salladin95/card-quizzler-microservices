package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/config"
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
	config config.AppCfg    // Application configuration
}

// IApp defines the interface for the main application.
type IApp interface {
	Start()
}

// NewApp creates a new instance of the application.
func NewApp(cfg config.AppCfg, rabbit *amqp.Connection) IApp {
	return &App{
		server: echo.New(), // Initialize Echo server
		rabbit: rabbit,
		config: cfg,
	}
}

// Start initializes and starts the application.
func (app *App) Start() {
	// Setup middlewares for the Echo server.
	app.setupMiddlewares()

	// Setup routes for the Echo server.
	app.setupRoutes()

	// Start the Echo server in a goroutine.
	go func() {
		serverAddr := fmt.Sprintf(":%s", app.config.BROKER_SERVICE_PORT)
		if err := app.server.Start(serverAddr); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Println("********* SHUTTING DOWN THE SERVER **********")
		}
	}()

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
