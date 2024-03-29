package server

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/handlers"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/repositories"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/subscribers"
	"github.com/Salladin95/rmqtools"
	"github.com/go-redis/redis"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
	"time"
)

// App represents the main application structure.
type App struct {
	broker rmqtools.MessageBroker
	config *config.Config
	db     *gorm.DB
	redis  *redis.Client
}

// IApp defines the interface for the main application.
type IApp interface {
	Start()
}

// NewApp creates a new instance of the application.
func NewApp(
	cfg *config.Config, // Application configuration
	rabbit *amqp091.Connection, // RabbitMQ connection
	db *gorm.DB, // GormDB client
	redisClient *redis.Client, // Redis client
) IApp {
	return &App{
		broker: rmqtools.NewMessageBroker(
			rabbit,
			constants.AmqpExchange,
			constants.AmqpQueue,
		),
		config: cfg,
		db:     db,
		redis:  redisClient,
	}
}

// Start initializes and starts the application.
func (app *App) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	repo := repositories.NewRepo(app.db, app.broker)

	go subscribers.NewMessageBrokerSubscribers(app.broker, repo).Listen(ctx)

	// Start the application by invoking the gRPC listener.
	app.gRPCListen(repo)
}

// gRPCListen sets up a gRPC server and listens for incoming requests on the specified port.
func (app *App) gRPCListen(repo repositories.Repository) {
	// Create a TCP listener for the specified gRPC port.
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", app.config.GrpcPort))
	if err != nil {
		msg := fmt.Sprintf(
			"failed to listen tcp port - %s. Err - %s",
			app.config.GrpcPort,
			err.Error(),
		)
		log.Fatalf(msg)   // Fatal log and exit if listener creation fails
		lib.LogError(msg) // Log error message
	}

	// Create a new gRPC server instance.
	gRPCServer := grpc.NewServer()

	// Register the quizzler server implementation with the gRPC server.
	handlers.RegisterQuizzlerServer(
		gRPCServer,
		repo,
		app.broker,
	)

	// Log a message indicating that the gRPC server has started.
	lib.LogInfo(fmt.Sprintf("gRPC Server started on port %s", app.config.GrpcPort))

	// Start serving gRPC requests on the listener.
	if err := gRPCServer.Serve(listener); err != nil {
		msg := fmt.Sprintf("Failed to listen for gRPC: %v", err)
		lib.LogError(msg) // Log error message
		log.Fatalf(msg)   // Fatal log and exit if server fails to serve
	}
}
