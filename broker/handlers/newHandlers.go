package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

const (
	AmqpExchange          = "broker"
	defaultAuthServiceURL = "localhost:8090" //local
	//defaultAuthServiceURL = "auth-service:8090" //from container
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

func (bh *brokerHandlers) GetGRPCClientConn() (*grpc.ClientConn, error) {
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
	return conn, nil
	//defer conn.Close()
	//return auth.NewAuthClient(conn), nil
}
