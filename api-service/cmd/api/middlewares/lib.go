package middlewares

import (
	"errors"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/goErrorHandler"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
	"time"
)

// GetJwtConfig returns a configuration for JWT authentication.
// It takes a key as a parameter for the signing key.
func GetJwtConfig(key string) echojwt.Config {
	return echojwt.Config{
		SigningKey: []byte(key), // Set the signing key for JWT
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(entities.JwtUserClaims) // Create new JWT claims using the JwtUserClaims struct
		},
		ErrorHandler: func(c echo.Context, err error) error {
			// Handle errors during JWT authentication and return an unauthorized error
			return goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
		},
	}
}

// validateTokenString validates a JWT token string using the provided secret key.
// It returns the decoded claims if the token is valid, otherwise, returns an error.
func validateTokenString(tokenString string, secret string) (*entities.JwtUserClaims, error) {
	// Get JWT configuration with the provided secret key
	config := GetJwtConfig(secret)

	// Parse and validate the token with custom claims and signing key
	token, err := jwt.ParseWithClaims(
		tokenString,
		&entities.JwtUserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return config.SigningKey, nil
		})

	if err != nil || !token.Valid {
		// Log and return an unauthorized error if the token is invalid
		log.Errorf("Token validation failed: %v - Token: %s", err, tokenString)
		return nil, goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
	}

	// Extract claims from the token
	claims, ok := token.Claims.(*entities.JwtUserClaims)
	if !ok {
		// Log and return an unauthorized error if claims are invalid
		log.Printf("Invalid claims in the token: %s\n", tokenString)
		return nil, goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
	}

	return claims, nil
}

// clearCookies removes refresh-token cookie from the provided echo.Context.
// It takes the echo.Context as parameter.
func clearCookies(c echo.Context) {
	// Reset the specified cookies
	ResetCookies(c, "refresh-token")
	log.Info("Cleaning cookies")
}

// ResetCookies clears cookies by setting new cookies with an expiration time in the past.
// It takes echo.Context and cookie names as parameters.
func ResetCookies(c echo.Context, cookieNames ...string) {
	for _, cookieName := range cookieNames {
		// Create a new cookie with an expiration time in the past
		cookie := new(http.Cookie)
		cookie.Name = cookieName
		cookie.Value = ""
		cookie.Expires = time.Now().Add(-24 * time.Hour) // Set expiration time to the past
		// Set the cookie in the response
		c.SetCookie(cookie)
	}
}

// ExtractAccessToken extracts the access token from the Authorization header in the echo.Context.
// It returns the extracted access token and an error, if any.
func ExtractAccessToken(c echo.Context) (string, error) {
	// Get the Authorization header
	authHeader := c.Request().Header.Get("Authorization")

	// Split the Authorization header to retrieve the token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		// If the token is missing or has an incorrect format, clear cookies and return an error
		clearCookies(c)
		log.Printf("Token is missing or has an incorrect format\n")
		return "", goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, errors.New("token is missing or has an incorrect format"))
	}

	// Extract and return the Bearer token
	bearerToken := strings.TrimPrefix(tokenParts[1], "Bearer ")
	bearerToken = strings.Trim(bearerToken, "\"")
	return bearerToken, nil
}

// ExtractRefreshToken extracts the refresh token from the "refresh-token" cookie in the echo.Context.
// It returns the extracted refresh token and an error, if any.
func ExtractRefreshToken(c echo.Context) (string, error) {
	// Get the "refresh-token" cookie from the request
	tokenCookie, err := c.Cookie("refresh-token")
	if err != nil || tokenCookie.Value == "" {
		// If the refresh token cookie is missing or empty, clear cookies and return an error
		log.Infof("Refresh token cookie is missing or empty\n")
		clearCookies(c)
		return "", goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
	}

	// Extract and return the refresh token
	tokenString := tokenCookie.Value
	return tokenString, nil
}

func generateValidatorLog(message string, level string, method string) entities.LogMessage {
	var logMessage entities.LogMessage
	return logMessage.GenerateLog(message, level, method, "token validator")
}
