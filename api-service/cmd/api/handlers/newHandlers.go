package handlers

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AmqpExchange is the name of the AMQP exchange used by the api-service.
const (
	AmqpExchange = "api-service"
)

// BrokerHandlersInterface defines the interface for api-service-related HTTP handlers.
type BrokerHandlersInterface interface {
	SignIn(c echo.Context) error
	SignUp(c echo.Context) error
	Refresh(c echo.Context) error
	GetUsers(c echo.Context) error
	GetUserById(c echo.Context) error
	GetProfile(c echo.Context) error
}

// brokerHandlers implements the BrokerHandlersInterface.
type brokerHandlers struct {
	rabbit       *amqp.Connection
	config       *config.Config
	cacheManager cacheManager.CacheManager
}

// NewHandlers creates a new instance of BrokerHandlersInterface.
func NewHandlers(
	cfg *config.Config,
	rabbit *amqp.Connection,
	cacheManager cacheManager.CacheManager,
) BrokerHandlersInterface {
	return &brokerHandlers{
		rabbit:       rabbit,
		config:       cfg,
		cacheManager: cacheManager,
	}
}

// GetGRPCClientConn establishes a gRPC client connection using the specified URL and returns the connection.
func (bh *brokerHandlers) GetGRPCClientConn() (*grpc.ClientConn, error) {
	// Dial a gRPC server using the provided URL and insecure transport credentials
	conn, err := grpc.Dial(
		bh.config.AppCfg.UserServiceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		// Handle the error and return an OperationFailure error using the goErrorHandler package
		return nil, goErrorHandler.OperationFailure("connect to gRPC", err)
	}

	return conn, nil
}