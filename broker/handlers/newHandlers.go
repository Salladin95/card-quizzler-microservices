package handlers

import (
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	AmqpExchange = "broker"
)

type BrokerHandlers interface {
	SignIn(c echo.Context) error
	SignUp(c echo.Context) error
}

type brokerHandlers struct {
	rabbit *amqp.Connection
}

func NewHandlers(rabbit *amqp.Connection) BrokerHandlers {
	return &brokerHandlers{rabbit: rabbit}
}
