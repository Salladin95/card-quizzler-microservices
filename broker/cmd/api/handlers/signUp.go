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

func (bh *brokerHandlers) SignUp(c echo.Context) error {
	fmt.Println("******* broker - start processing signUp request ********")
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
	ah := auth.NewAuthClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Ensure the context is canceled when done.
	res, err := ah.SignUp(ctx, &auth.SignUpRequest{
		Payload: signUpDTO.ToAuthPayload(),
	})
	if err != nil {
		return goErrorHandler.OperationFailure("sign up", err)
	}
	return buildUserResponse(c, res)
}
