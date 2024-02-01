package handlers

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/config"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	AmqpExchange = "broker"
)

type BrokerHandlersInterface interface {
	SignIn(c echo.Context) error
	SignUp(c echo.Context) error
}

type brokerHandlers struct {
	rabbit *amqp091.Connection
	config config.AppCfg
}

func NewHandlers(cfg config.AppCfg, rabbit *amqp091.Connection) BrokerHandlersInterface {
	return &brokerHandlers{
		rabbit: rabbit,
		config: cfg,
	}
}

func (bh *brokerHandlers) GetGRPCClientConn() (*grpc.ClientConn, error) {
	fmt.Printf("******* broker; connecting to gRPC on url - %s ********\n", bh.config.AUTH_SERVICE_URL)
	conn, err := grpc.Dial(bh.config.AUTH_SERVICE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	fmt.Printf("******* broker; connected to gRPC on url - %s !!!!!!\n", bh.config.AUTH_SERVICE_URL)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get auth client error - %v", err)
	}
	return conn, nil
}
