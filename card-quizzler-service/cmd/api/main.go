package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/handlers"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/repositories"
	"github.com/Salladin95/rmqtools"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	server := echo.New()
	cfg, err := config.NewConfig()

	if err != nil {
		server.Logger.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	db, err := gorm.Open(postgres.Open(cfg.DbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	repo := repositories.NewRepo(db)

	// Connect to RabbitMQ server using the provided URL.
	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.RabbitUrl)
	if err != nil {
		log.Println(err) // Log error if RabbitMQ connection fails
		os.Exit(1)       // Exit program if RabbitMQ connection fails
	}

	// Create a new gRPC server instance.
	gRPCServer := grpc.NewServer()

	broker := rmqtools.NewMessageBroker(rabbitConn, constants.AmqpExchange, constants.AmqpQueue)

	// Register the AuthServer implementation with the gRPC server.
	handlers.RegisterQuizzlerServer(gRPCServer, repo, broker)

	//migrations.Migrate(db)

	// Start the Echo server in a goroutine.
	go func() {
		serverAddr := fmt.Sprintf(":%s", cfg.GrpcPort)
		if err := server.Start(serverAddr); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Panic("shutting down the server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Gracefully shut down the Echo server.
	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatal(err)
	}
}
