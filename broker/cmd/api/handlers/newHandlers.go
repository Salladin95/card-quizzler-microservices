package handlers

import (
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/config"
	"github.com/Salladin95/goErrorHandler"
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
	conn, err := grpc.Dial(bh.config.AUTH_SERVICE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, goErrorHandler.OperationFailure("connect to gRPC", err)
	}

	if err != nil {
		return nil, goErrorHandler.OperationFailure("get auth client", err)
	}
	return conn, nil
}
