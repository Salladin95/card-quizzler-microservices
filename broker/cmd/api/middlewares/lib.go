package middlewares

import (
	"errors"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/entities"
	"github.com/Salladin95/goErrorHandler"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
	"time"
)

func GetJwtConfig(key string) echojwt.Config {
	return echojwt.Config{
		SigningKey: []byte(key),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(entities.JwtUserClaims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
		},
	}
}

func validateTokenString(tokenString string, secret string) (*entities.JwtUserClaims, error) {
	config := GetJwtConfig(secret)
	token, err := jwt.ParseWithClaims(tokenString, &entities.JwtUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.SigningKey, nil
	})

	if err != nil || !token.Valid {
		log.Errorf("token - %s has failed validation: %v", err, tokenString)
		return nil, goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
	}

	// Generate claims
	claims, ok := token.Claims.(*entities.JwtUserClaims)
	if !ok {
		log.Printf("********** invalid claims *************\n")
		return nil, goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
	}
	return claims, nil
}

func clearSessionAndCookies(c echo.Context) {
	ResetCookies(c, "refresh-token")
	log.Info("cleaning session & cookies")
}

// ResetCookies clears cookies by setting new cookies with an expiration time in the past.
func ResetCookies(c echo.Context, cookieNames ...string) {
	for _, cookieName := range cookieNames {
		// Set a new cookie with an expiration time in the past
		cookie := new(http.Cookie)
		cookie.Name = cookieName
		cookie.Value = ""
		cookie.Expires = time.Now().Add(-24 * time.Hour)
		c.SetCookie(cookie)
	}
}

func ExtractAccessToken(c echo.Context) (string, error) {
	// Get the Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		clearSessionAndCookies(c)
		log.Printf("********** token is empty *************\n")
		return "", goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, errors.New("token is missing"))
	}
	bearerToken := strings.TrimPrefix(tokenParts[1], "Bearer ")
	bearerToken = strings.Trim(bearerToken, "\"")
	return bearerToken, nil
}

func ExtractRefreshToken(c echo.Context) (string, error) {
	tokenCookie, err := c.Cookie("refresh-token")
	if err != nil || tokenCookie.Value == "" {
		log.Infof("****** refresh token cookie is missing **********\n")
		clearSessionAndCookies(c)
		return "", goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
	}
	tokenString := tokenCookie.Value
	return tokenString, err
}
