package handlers

import (
	"errors"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	userService "github.com/Salladin95/card-quizzler-microservices/api-service/user"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (ah *apiHandlers) GetUserById(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "getUserById")

	uid := c.Param("id")

	user, err := ah.cacheManager.GetUserById(ctx, uid)

	if err == nil {
		return handleCacheResponse(c, user)
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn()
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Make a gRPC call to the SignIn method of the Auth service
	res, err := userService.NewUserServiceClient(clientConn).GetUserById(ctx, &userService.RequestWithID{
		Id: uid,
	})
	if err != nil {
		return goErrorHandler.OperationFailure("GetUserById", err)
	}
	var unmarshalTo []entities.UserResponse
	return handleGRPCResponse(c, res, unmarshalTo)
}

func (ah *apiHandlers) GetProfile(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "GetProfile")

	// Retrieve user claims from the context
	claims, ok := c.Get("user").(*entities.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(
			goErrorHandler.ErrUnauthorized,
			errors.New("refresh, failed to cast claims"),
		)
	}

	user, err := ah.cacheManager.GetUserById(ctx, claims.Id)

	if err == nil {
		return handleCacheResponse(c, user)
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn()
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Make a gRPC call to the SignIn method of the Auth service
	res, err := userService.NewUserServiceClient(clientConn).GetUserById(ctx, &userService.RequestWithID{
		Id: claims.Id,
	})
	if err != nil {
		return goErrorHandler.OperationFailure("GetProfile", err)
	}
	var unmarshalTo []entities.UserResponse
	return handleGRPCResponse(c, res, unmarshalTo)
}

func (ah *apiHandlers) UpdateEmail(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "UpdateEmail")

	uid := c.Param("id")

	var dto entities.UpdateEmailDto
	err := lib.BindBodyAndVerify(c, &dto)
	if err != nil {
		return err
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn()
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Make a gRPC call to the SignIn method of the Auth service
	res, err := userService.NewUserServiceClient(clientConn).UpdateEmail(ctx, &userService.UpdateEmailRequest{
		Payload: dto.ToPayload(uid),
	})
	if err != nil {
		return goErrorHandler.OperationFailure("UpdateEmail", err)
	}
	var unmarshalTo []entities.UserResponse
	return handleGRPCResponse(c, res, unmarshalTo)
}

func (ah *apiHandlers) UpdatePassword(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "UpdatePassword")

	uid := c.Param("id")

	var dto entities.UpdatePasswordDto
	err := lib.BindBodyAndVerify(c, &dto)
	if err != nil {
		return err
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn()
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Make a gRPC call to the SignIn method of the Auth service
	res, err := userService.NewUserServiceClient(clientConn).UpdatePassword(ctx, &userService.UpdatePasswordRequest{
		Payload: dto.ToPayload(uid),
	})
	if err != nil {
		return goErrorHandler.OperationFailure("UpdatePassword", err)
	}
	var unmarshalTo []entities.UserResponse
	return handleGRPCResponse(c, res, unmarshalTo)
}

func (ah *apiHandlers) ResetPassword(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "ResetPassword")

	var dto entities.ResetPasswordDto
	err := lib.BindBodyAndVerify(c, &dto)
	if err != nil {
		return err
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn()
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Make a gRPC call to the SignIn method of the Auth service
	res, err := userService.NewUserServiceClient(clientConn).ResetPassword(ctx, &userService.ResetPasswordRequest{
		Payload: dto.ToPayload(),
	})
	if err != nil {
		return goErrorHandler.OperationFailure("ResetPassword", err)
	}
	var unmarshalTo []entities.UserResponse
	return handleGRPCResponse(c, res, unmarshalTo)
}

func (ah *apiHandlers) RequestEmailVerification(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "RequestEmailVerification")
	var dto entities.RequestEmailVerificationDto
	err := lib.BindBodyAndVerify(c, &dto)
	if err != nil {
		return err
	}

	err = ah.broker.PushToQueue(ctx, constants.RequestEmailVerificationCommand, dto)

	if err != nil {
		return err
	}
	ah.log(ctx, "generated event for email verification", "info", "RequestEmailVerification")
	return c.String(http.StatusNoContent, "Verification code is sent.")
}
