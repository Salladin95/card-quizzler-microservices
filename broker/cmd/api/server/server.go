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

type App struct {
	server *echo.Echo
	rabbit *amqp.Connection
	config config.AppCfg
}

type IApp interface {
	Start()
}

func NewApp(cfg config.AppCfg, rabbit *amqp.Connection) IApp {
	return &App{
		server: echo.New(),
		rabbit: rabbit,
		config: cfg,
	}
}

func (app *App) Start() {
	app.setupMiddlewares()
	app.setupRoutes()

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.server.Shutdown(ctx); err != nil {
		app.server.Logger.Fatal(err)
	}
}
