package handlers

import (
	"errors"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Refresh handles the HTTP request for token refresh.
func (bh *brokerHandlers) Refresh(c echo.Context) error {
	fmt.Println("******* api-service - start processing refresh request ***************")

	// Retrieve user claims from the context
	claims, ok := c.Get("user").(*entities.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, errors.New("refresh, failed to cast claims"))
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
