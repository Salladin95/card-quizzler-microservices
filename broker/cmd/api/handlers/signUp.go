package handlers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/auth"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"time"
)

// SignUp handles the HTTP request for user sign-up.
func (bh *brokerHandlers) SignUp(c echo.Context) error {
	fmt.Println("******* broker - start processing signUp request ********")

	// Parse the request body into SignUpDto
	var signUpDTO entities.SignUpDto
	if err := lib.BindBodyAndVerify(c, &signUpDTO); err != nil {
		return err
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

	// Make a gRPC call to the SignUp method of the Auth service
	res, err := auth.NewAuthClient(clientConn).SignUp(ctx, &auth.SignUpRequest{
		Payload: signUpDTO.ToAuthPayload(),
	})
	if err != nil {
		return goErrorHandler.OperationFailure("sign up", err)
	}

	// Check the response code from the Auth service
	resCode := int(res.GetCode())
	if resCode >= 400 {
		return c.JSON(resCode, entities.JsonResponse{Message: res.GetMessage()})
	}

	// Unmarshal the user response data from the gRPC response
	var userResponse entities.UserResponse
	err = lib.UnmarshalData(res.GetData(), &userResponse)
	if err != nil {
		return goErrorHandler.OperationFailure("unmarshal data", err)
	}

	// Respond with the JSON data containing the user information
	return c.JSON(resCode, entities.JsonResponse{Message: res.GetMessage(), Data: userResponse})
}
