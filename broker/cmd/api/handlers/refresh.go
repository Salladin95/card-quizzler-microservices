package handlers

import (
	"errors"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/entities"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (bh *brokerHandlers) Refresh(c echo.Context) error {
	fmt.Println("******* broker - start processing refresh request ***************")
	claims, ok := c.Get("user").(*entities.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, errors.New("refresh, failed to cast claims"))
	}
	tokens, err := GenerateTokenPair(claims.Name, claims.Email, claims.Id, bh.config.JwtCfg)
	if err != nil {
		return err
	}

	SetHttpOnlyCookie(
		c,
		"refresh-token",
		tokens.RefreshToken,
		bh.config.JwtCfg.RefreshTokenExpTime,
		"/",
	)

	return c.JSON(http.StatusOK, entities.SignInResponse{
		AccessToken: tokens.AccessToken,
	})
}
