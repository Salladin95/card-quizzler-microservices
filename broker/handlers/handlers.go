package handlers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/auth"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	SignInKey = "auth.sign-in.command"
	SignUpKey = "auth.sign-up.command"
)

type SignInDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"min=6,required"`
}

type SighUpDto struct {
	Name     string `json:"name"  validate:"required,min=1"`
	Password string `json:"password"  validate:"required,min=6"`
	Email    string `json:"email"  validate:"required,email"`
	Birthday string `json:"birthday"  validate:"required,min=1"`
}

type JsonResponse struct {
	message string
}

func (bh *brokerHandlers) SignIn(c echo.Context) error {
	fmt.Println("******* broker - start processing signIn request ***************")
	var signInDTO SignInDto

	// Read the request body and unmarshal it into the corresponding DTO
	if err := c.Bind(&signInDTO); err != nil {
		return fmt.Errorf("Error binding request body: %v\n", err)
	}

	clientConn, err := bh.GetGRPCClientConn()
	defer clientConn.Close()
	if err != nil {
		return err
	}
	ah := auth.NewAuthClient(clientConn)

	// Use a longer timeout for the gRPC call, adjust as needed
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := ah.SignIn(ctx, &auth.SignInRequest{
		Payload: &auth.SignInPayload{
			Email:    signInDTO.Email,
			Password: signInDTO.Password,
		},
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, JsonResponse{message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (bh *brokerHandlers) SignUp(c echo.Context) error {
	fmt.Println("******* broker - start processing signUp request ********")
	var signUpDTO SighUpDto

	// Read the request body and unmarshal it into the corresponding DTO
	if err := c.Bind(&signUpDTO); err != nil {
		return fmt.Errorf("Error binding request body: %v\n", err)
	}

	clientConn, err := bh.GetGRPCClientConn()
	defer clientConn.Close()
	if err != nil {
		return err
	}
	ah := auth.NewAuthClient(clientConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := ah.SignUp(ctx, &auth.SignUpRequest{
		Payload: &auth.SignUpPayload{
			Email:    signUpDTO.Email,
			Password: signUpDTO.Password,
			Name:     signUpDTO.Name,
			Birthday: signUpDTO.Birthday,
		},
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, JsonResponse{message: err.Error()})
	}

	return c.JSON(http.StatusOK, JsonResponse{message: res.GetMessage()})
}
