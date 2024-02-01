package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Salladin95/goErrorHandler"
	"github.com/Salladin95/rmqtools"
	"github.com/labstack/echo/v4"
	"net/http"
)

// pushToQueue pushes data to a named queue in the AMQP exchange.
func (bh *brokerHandlers) pushToQueue(ctx context.Context, name string, data []byte) error {
	// Create a new EventEmitter for the specified AMQP exchange.
	emitter, err := rmqtools.NewEventEmitter(bh.rabbit, AmqpExchange)
	if err != nil {
		return goErrorHandler.OperationFailure("create event emitter", err)
	}

	// Push the data to the named queue using the EventEmitter.
	err = emitter.Push(ctx, name, data)
	if err != nil {
		return goErrorHandler.OperationFailure("push event", err)
	}

	return nil
}

// pushToQueueFromEndpoint handles HTTP requests to push data to a named queue from an endpoint.
func (bh *brokerHandlers) pushToQueueFromEndpoint(c echo.Context, key string) error {
	var requestDTO interface{}

	// Read the request body and unmarshal it into the corresponding DTO.
	if err := c.Bind(&requestDTO); err != nil {
		return goErrorHandler.BindRequestToBodyFailure(err)
	}

	// Marshal the DTO struct into JSON.
	marshalledDto, err := json.Marshal(requestDTO)
	if err != nil {
		return goErrorHandler.OperationFailure("marshal dto", err)
	}

	// Push the marshalled DTO to the named queue using pushToQueue method.
	err = bh.pushToQueue(c.Request().Context(), key, marshalledDto)
	if err != nil {
		return err
	}

	// Respond with a success message indicating the push operation.
	c.String(http.StatusOK, fmt.Sprintf("PUSHED TO QUEUE FROM %s", key))
	return nil
}
