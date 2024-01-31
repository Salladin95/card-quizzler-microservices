package handlers

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/auth"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

const (
	AmqpExchange          = "broker"
	defaultAuthServiceURL = "auth-service:8090"
)

type BrokerHandlersInterface interface {
	SignIn(c echo.Context) error
	SignUp(c echo.Context) error
}

type brokerHandlers struct {
	rabbit *amqp091.Connection
}

func NewHandlers(rabbit *amqp091.Connection) BrokerHandlersInterface {
	return &brokerHandlers{
		rabbit: rabbit,
	}
}

func (bh *brokerHandlers) GetGRPCClientConn() (auth.AuthClient, error) {
	authUrl := os.Getenv("AUTH_SERVICE_URL")
	if authUrl == "" {
		authUrl = defaultAuthServiceURL
	}
	fmt.Printf("******* broker; connecting to gRPC on url - %s ********\n", authUrl)
	conn, err := grpc.Dial(authUrl, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	fmt.Printf("******* broker; connected to gRPC on url - %s !!!!!!\n", authUrl)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get auth client error - %v", err)
	}
	defer conn.Close()

	return auth.NewAuthClient(conn), nil
}
