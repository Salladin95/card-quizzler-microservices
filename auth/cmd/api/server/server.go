package server

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/handlers"
	auth "github.com/Salladin95/card-quizzler-microservices/auth-service/proto"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"log"
	"net"
)

// App represents the main application structure.
type App struct {
	rabbit *amqp091.Connection // RabbitMQ connection instance
	config config.AppCfg       // Application configuration
}

// IApp defines the interface for the main application.
type IApp interface {
	Start()
}

// NewApp creates a new instance of the application.
func NewApp(cfg config.AppCfg, rabbit *amqp091.Connection) IApp {
	return &App{
		rabbit: rabbit,
		config: cfg,
	}
}

// Start initializes and starts the application.
func (app *App) Start() {
	// Start the application by invoking the gRPC listener.
	app.gRPCListen()
}

// gRPCListen sets up a gRPC server and listens for incoming requests on the specified port.
func (app *App) gRPCListen() {
	// Create a TCP listener for the specified gRPC port.
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", app.config.GRPC_PORT))
	if err != nil {
		log.Fatalf("failed to listen tcp port - %s. Err - %s", app.config.GRPC_PORT, err.Error())
	}

	// Create a new gRPC server instance.
	gRPCServer := grpc.NewServer()

	// Register the AuthServer implementation with the gRPC server.
	auth.RegisterAuthServer(gRPCServer, &handlers.AuthServer{})

	// Log a message indicating that the gRPC server has started.
	log.Printf("gRPC Server started on port %s", app.config.GRPC_PORT)

	// Start serving gRPC requests on the listener.
	if err := gRPCServer.Serve(listener); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
