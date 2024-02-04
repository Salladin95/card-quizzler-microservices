package server

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/handlers"
	user "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/repository"
	auth "github.com/Salladin95/card-quizzler-microservices/auth-service/proto"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

// App represents the main application structure.
type App struct {
	rabbit *amqp091.Connection
	config *config.Config
	db     *mongo.Client
}

// IApp defines the interface for the main application.
type IApp interface {
	Start()
}

// NewApp creates a new instance of the application.
func NewApp(
	cfg *config.Config,
	rabbit *amqp091.Connection,
	db *mongo.Client,
) IApp {
	return &App{
		rabbit: rabbit,
		config: cfg,
		db:     db,
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
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", app.config.AppCfg.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen tcp port - %s. Err - %s", app.config.AppCfg.GrpcPort, err.Error())
	}

	// Create a new gRPC server instance.
	gRPCServer := grpc.NewServer()

	// Register the AuthServer implementation with the gRPC server.
	auth.RegisterAuthServer(gRPCServer, &handlers.AuthServer{Repo: user.NewUserRepository(app.db, app.config.MongoCfg)})

	// Log a message indicating that the gRPC server has started.
	log.Printf("gRPC Server started on port %s", app.config.AppCfg.GrpcPort)

	// Start serving gRPC requests on the listener.
	if err := gRPCServer.Serve(listener); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
