package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/auth"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (bh *brokerHandlers) SignIn(c echo.Context) error {
	fmt.Println("******* broker - start processing signIn request ***************")
	var signInDTO entities.SignInDto
	if err := lib.BindBodyAndVerify(c, &signInDTO); err != nil {
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
	res, err := ah.SignIn(ctx, &auth.SignInRequest{
		Payload: signInDTO.ToAuthPayload(),
	})
	if err != nil {
		return goErrorHandler.OperationFailure("sign in", err)
	}

	resCode := int(res.GetCode())
	if resCode >= 400 {
		return c.JSON(resCode, entities.JsonResponse{Message: res.GetMessage()})
	}
	var signedInUser entities.UserResponse
	err = json.Unmarshal(res.GetData(), &signedInUser)
	if err != nil {
		return goErrorHandler.OperationFailure("unmarshal data", err)
	}
	tokens, err := GenerateTokenPair(signedInUser.Name, signedInUser.Email, signedInUser.ID, bh.config.JwtCfg)
	SetHttpOnlyCookie(
		c,
		"refresh-token",
		tokens.RefreshToken,
		//bh.config.JwtCfg.RefreshTokenExpTime,
		1*time.Minute,
		"/",
	)

	return c.JSON(http.StatusOK, entities.SignInResponse{
		AccessToken: tokens.AccessToken,
	})
}
