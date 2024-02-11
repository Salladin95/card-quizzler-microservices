package messageBroker

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
	"github.com/Salladin95/rmqtools"
	"github.com/rabbitmq/amqp091-go"
)

type messageBroker struct {
	rabbitConn *amqp091.Connection
}

type MessageBroker interface {
	PushToQueue(ctx context.Context, name string, data interface{}) error
	GenerateLogEvent(ctx context.Context, logMessage entities.LogMessage) error
	ListenForUpdates(topics []string, mh func(routingKey string, payload []byte))
}

func NewMessageBroker(rabbitConn *amqp091.Connection) MessageBroker {
	return &messageBroker{rabbitConn: rabbitConn}
}

// PushToQueue pushes data to a named queue in RabbitMQ using an EventEmitter.
// It takes a context, the routingKey, and the data to be pushed.
// It returns an error if any occurs.
func (mb *messageBroker) PushToQueue(ctx context.Context, routingKey string, data interface{}) error {
	// Create a new EventEmitter for the specified AMQP exchange.
	emitter, err := rmqtools.NewEventEmitter(mb.rabbitConn, constants.AmqpExchange)
	if err != nil {
		return goErrorHandler.OperationFailure("create event emitter", err)
	}

	// Marshal the data into JSON format.
	mData, err := lib.MarshalData(data)
	if err != nil {
		return err
	}

	// Push the data to the named queue using the EventEmitter.
	err = emitter.Push(ctx, routingKey, mData)
	if err != nil {
		return goErrorHandler.OperationFailure("push event", err)
	}

	return nil
}
