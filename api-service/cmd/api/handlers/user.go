package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	userService "github.com/Salladin95/card-quizzler-microservices/api-service/user"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"time"
)

func (bh *brokerHandlers) GetUsers(c echo.Context) error {
	fmt.Println("******* api-service - start processing GetUsers request ***************")

	users, err := bh.cacheManager.GetUsers()

	if err == nil {
		return handleCacheResponse(c, users)
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from brokerHandlers.
	clientConn, err := bh.GetGRPCClientConn()
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Create a context with a timeout for the gRPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Ensure the context is canceled when done.

	// Make a gRPC call to the SignIn method of the Auth service
	res, err := userService.NewUserServiceClient(clientConn).GetUsers(ctx, &userService.EmptyRequest{})
	if err != nil {
		return goErrorHandler.OperationFailure("GetUsers", err)
	}
	var unmarshalTo []*entities.UserResponse
	return handleGRPCResponse(c, res, unmarshalTo)
}

func (bh *brokerHandlers) GetUserById(c echo.Context) error {
	fmt.Println("******* api-service - start processing GetUserById request ***************")

	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		return goErrorHandler.OperationFailure("parse id", err)
	}

	user, err := bh.cacheManager.GetUserById(uid.String())

	if err == nil {
		return handleCacheResponse(c, user)
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from brokerHandlers.
	clientConn, err := bh.GetGRPCClientConn()
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Create a context with a timeout for the gRPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Ensure the context is canceled when done.

	// Make a gRPC call to the SignIn method of the Auth service
	res, err := userService.NewUserServiceClient(clientConn).GetUserById(ctx, &userService.ID{
		Id: uid.String(),
	})
	if err != nil {
		return goErrorHandler.OperationFailure("GetUserById", err)
	}
	var unmarshalTo []entities.UserResponse
	return handleGRPCResponse(c, res, unmarshalTo)
}

func (bh *brokerHandlers) GetProfile(c echo.Context) error {
	fmt.Println("******* api-service - start processing GetProfile request ***************")

	// Retrieve user claims from the context
	claims, ok := c.Get("user").(*entities.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, errors.New("refresh, failed to cast claims"))
	}

	user, err := bh.cacheManager.GetUserById(claims.Id.String())

	if err == nil {
		return handleCacheResponse(c, user)
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from brokerHandlers.
	clientConn, err := bh.GetGRPCClientConn()
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Create a context with a timeout for the gRPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Ensure the context is canceled when done.

	// Make a gRPC call to the SignIn method of the Auth service
	res, err := userService.NewUserServiceClient(clientConn).GetUserById(ctx, &userService.ID{
		Id: claims.Id.String(),
	})
	if err != nil {
		return goErrorHandler.OperationFailure("GetUserById", err)
	}
	var unmarshalTo []entities.UserResponse
	return handleGRPCResponse(c, res, unmarshalTo)
}
