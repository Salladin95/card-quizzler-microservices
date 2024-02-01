package handlers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/auth"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/lib"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (bh *brokerHandlers) SignIn(c echo.Context) error {
	fmt.Println("******* broker - start processing signIn request ***************")
	var signInDTO entities.SignInDto

	// Read the request body and unmarshal it into the corresponding DTO (Data Transfer Object).
	// The BindBodyAndVerify function binds the request body to the provided DataWithVerify interface,
	// in this case, the signInDTO structure, and then calls the Verify method for additional validation.
	// If any error occurs during binding or verification, it is returned, and the function exits with an error response.
	if err := lib.BindBodyAndVerify(c, &signInDTO); err != nil {
		return err
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from brokerHandlers.
	clientConn, err := bh.GetGRPCClientConn()
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Create an instance of the Auth service client using the obtained gRPC client connection.
	ah := auth.NewAuthClient(clientConn)

	// Create a context with a timeout for the gRPC request.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Ensure the context is canceled when done.

	// Make a gRPC request to the SignIn method of the Auth service with the provided SignInRequest.
	res, err := ah.SignIn(ctx, &auth.SignInRequest{
		Payload: &auth.SignInPayload{
			Email:    signInDTO.Email,
			Password: signInDTO.Password,
		},
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (bh *brokerHandlers) SignUp(c echo.Context) error {
	fmt.Println("******* broker - start processing signUp request ********")
	var signUpDTO entities.SignUpDto

	// Read the request body and unmarshal it into the corresponding DTO (Data Transfer Object).
	// The BindBodyAndVerify function binds the request body to the provided DataWithVerify interface,
	// in this case, the signInDTO structure, and then calls the Verify method for additional validation.
	// If any error occurs during binding or verification, it is returned, and the function exits with an error response.
	if err := lib.BindBodyAndVerify(c, &signUpDTO); err != nil {
		return err
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from brokerHandlers.
	clientConn, err := bh.GetGRPCClientConn()
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Create an instance of the Auth service client using the obtained gRPC client connection.
	ah := auth.NewAuthClient(clientConn)

	// Create a context with a timeout for the gRPC request.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Ensure the context is canceled when done.

	// Make a gRPC request to the SignUp method of the Auth service with the provided SignUpRequest.
	res, err := ah.SignUp(ctx, &auth.SignUpRequest{
		Payload: &auth.SignUpPayload{
			Email:    signUpDTO.Email,
			Password: signUpDTO.Password,
			Name:     signUpDTO.Name,
			Birthday: signUpDTO.Birthday,
		},
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, entities.JsonResponse{Message: res.GetMessage()})
}
