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

func (bh *brokerHandlers) pushToQueue(ctx context.Context, name string, data []byte) error {
	emitter, err := rmqtools.NewEventEmitter(bh.rabbit, AmqpExchange)
	if err != nil {
		return goErrorHandler.OperationFailure("create event emitter", err)

	}
	err = emitter.Push(ctx, name, data)
	if err != nil {
		return goErrorHandler.OperationFailure("push event", err)
	}
	return nil
}

func (bh *brokerHandlers) pushToQueueFromEndpoint(c echo.Context, key string) error {
	var requestDTO interface{}

	// Read the request body and unmarshal it into the corresponding DTO
	if err := c.Bind(&requestDTO); err != nil {
		return goErrorHandler.BindRequestToBodyFailure(err)
	}

	// Marshal the DTO struct into JSON
	marshalledDto, err := json.Marshal(requestDTO)
	if err != nil {
		return goErrorHandler.OperationFailure("marshal dto", err)
	}

	err = bh.pushToQueue(c.Request().Context(), key, marshalledDto)
	if err != nil {
		return err
	}

	c.String(http.StatusOK, fmt.Sprintf("PUSHED TO QUEUE FROM %s", key))
	return nil
}
