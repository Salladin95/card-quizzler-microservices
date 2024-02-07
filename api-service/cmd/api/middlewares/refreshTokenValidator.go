package middlewares

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// RefreshTokenValidator returns an Echo middleware that validates the refresh token.
func RefreshTokenValidator(cacheManager cacheManager.CacheManager, secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		log.Info("RefreshTokenValidator - start token validation")
		return func(c echo.Context) error {
			// Extract the refresh token from the request
			refreshToken, err := ExtractRefreshToken(c)
			if err != nil {
				return err
			}

			// Validate the refresh token
			claims, err := validateTokenString(refreshToken, secret)
			if err != nil {
				clearCookies(c)
				cacheManager.ClearUserData(claims.Id.String())
				return err
			}

			// Retrieve cached refresh token
			cachedRefreshToken, err := cacheManager.RefreshToken(claims.Id.String())

			// Compare tokens
			if err != nil || cachedRefreshToken != refreshToken {
				log.Infof("********* Received token and token from cache don't match ***********")
				clearCookies(c)
				cacheManager.ClearUserData(claims.Id.String())
				return goErrorHandler.NewError(
					goErrorHandler.ErrUnauthorized,
					fmt.Errorf("received token and token from cache don't match"),
				)
			}

			log.Errorf("****** Refresh token has passed validation *******\n")

			// Set user claims in the context for the next handler
			c.Set("user", claims)
			return next(c)
		}
	}
}
