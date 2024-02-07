package handlers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	userService "github.com/Salladin95/card-quizzler-microservices/api-service/user"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// SignIn handles the HTTP request for user sign-in.
func (bh *brokerHandlers) SignIn(c echo.Context) error {
	fmt.Println("******* api-service - start processing signIn request ***************")

	// Parse the request body into SignInDto
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

	// Create a context with a timeout for the gRPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Ensure the context is canceled when done.

	// Make a gRPC call to the SignIn method of the Auth service
	res, err := userService.NewUserServiceClient(clientConn).SignIn(ctx, &userService.SignInRequest{
		Payload: signInDTO.ToAuthPayload(),
	})
	if err != nil {
		return goErrorHandler.OperationFailure("sign in", err)
	}

	// Check the response code from the Auth service
	resCode := int(res.GetCode())
	if resCode >= 400 {
		return c.JSON(resCode, entities.JsonResponse{Message: res.GetMessage()})
	}

	// Unmarshal the user response data from the gRPC response
	var signedInUser entities.UserResponse
	err = lib.UnmarshalData(res.GetData(), &signedInUser)
	if err != nil {
		return err
	}

	// Generate a token pair for the signed-in user
	tokens, err := GenerateTokenPair(signedInUser.Name, signedInUser.Email, signedInUser.ID, bh.config.JwtCfg)
	if err != nil {
		return err
	}

	// Set the refresh token as an HTTP-only cookie
	SetHttpOnlyCookie(
		c,
		"refresh-token",
		tokens.RefreshToken,
		bh.config.JwtCfg.RefreshTokenExpTime,
		"/",
	)

	// Set the access and refresh tokens in the cache
	err = bh.cacheManager.SetTokenPair(signedInUser.ID.String(), tokens)
	if err != nil {
		return err
	}

	// Respond with the access token in the JSON body
	return c.JSON(http.StatusOK, entities.SignInResponse{
		AccessToken: tokens.AccessToken,
	})
}
