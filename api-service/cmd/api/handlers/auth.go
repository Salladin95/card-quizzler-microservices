package handlers

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	userService "github.com/Salladin95/card-quizzler-microservices/api-service/user"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"net/http"
)

// SignIn handles the HTTP request for user sign-in.
func (bh *brokerHandlers) SignIn(c echo.Context) error {
	ctx := c.Request().Context()
	bh.log(ctx, "start processing request", "info", "signIn")

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

// SignUp handles the HTTP request for user sign-up.
func (bh *brokerHandlers) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	bh.log(ctx, "start processing request", "info", "signUp")

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

	// Make a gRPC call to the SignUp method of the Auth service
	res, err := userService.NewUserServiceClient(clientConn).SignUp(ctx, &userService.SignUpRequest{
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

// Refresh handles the HTTP request for token refresh.
func (bh *brokerHandlers) Refresh(c echo.Context) error {
	ctx := c.Request().Context()
	bh.log(ctx, "start processing request", "info", "refresh")

	// Extract the refresh token from the request
	refreshToken, err := lib.ExtractRefreshToken(c)
	if err != nil {
		return err
	}

	// Validate the refresh token
	claims, err := lib.ValidateTokenString(refreshToken, bh.config.JwtCfg.JWTRefreshSecret)
	if err != nil {
		lib.ClearCookies(c)
		bh.cacheManager.ClearUserRelatedCache(claims.Id.String())
		bh.log(
			ctx,
			fmt.Sprintf("clearing cookies and session. Err - %s", err.Error()),
			"error",
			"refresh",
		)
		return err
	}

	// Retrieve cached refresh token
	cachedRefreshToken, err := bh.cacheManager.RefreshToken(claims.Id.String())

	// Compare tokens
	if err != nil || cachedRefreshToken != refreshToken {
		bh.log(
			ctx,

			"Received token and token from cache don't match. Clearing cache & session",
			"error",
			"refresh",
		)
		lib.ClearCookies(c)
		bh.cacheManager.ClearUserRelatedCache(claims.Id.String())
		return goErrorHandler.NewError(
			goErrorHandler.ErrUnauthorized,
			fmt.Errorf("received token and token from cache don't match"),
		)
	}

	// Generate a new pair of access and refresh tokens
	tokens, err := GenerateTokenPair(claims.Name, claims.Email, claims.Id, bh.config.JwtCfg)
	if err != nil {
		return err
	}

	// Set the new refresh token as an HTTP-only cookie
	SetHttpOnlyCookie(
		c,
		"refresh-token",
		tokens.RefreshToken,
		bh.config.JwtCfg.RefreshTokenExpTime,
		"/",
	)

	// Set the new token pair
	err = bh.cacheManager.SetTokenPair(claims.Id.String(), tokens)
	if err != nil {
		return err
	}

	// Respond with the JSON data containing the new access token
	return c.JSON(http.StatusOK, entities.SignInResponse{
		AccessToken: tokens.AccessToken,
	})
}
