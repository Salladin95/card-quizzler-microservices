package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Salladin95/rmqtools"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	serverPort       = 80
	defaultRabbitURL = "amqp://khalid:12345@localhost:5672"
	//defaultRabbitURL = "amqp://khalid:12345@rabbitmq:5672" // from container
)

type App struct {
	server *echo.Echo
	rabbit *amqp.Connection
}

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = defaultRabbitURL
	}

	rabbitConn, err := rmqtools.ConnectToRabbit(rabbitURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := App{
		server: echo.New(),
		rabbit: rabbitConn,
	}
	app.setupRoutes()

	go func() {
		if err := app.startServer(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Println("********* SHUTTING DOWN THE SERVER **********")
		}
	}()
	app.waitForInterruptSignal()
}

func (app *App) startServer() error {
	serverAddr := fmt.Sprintf(":%d", serverPort)
	return app.server.Start(serverAddr)
}

func (app *App) waitForInterruptSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
