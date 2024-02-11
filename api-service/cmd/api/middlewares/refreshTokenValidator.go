package middlewares

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/messageBroker"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"time"
)

// RefreshTokenValidator returns an Echo middleware that validates the refresh token.
func RefreshTokenValidator(broker messageBroker.MessageBroker, cacheManager cacheManager.CacheManager, secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		broker.GenerateLogEvent(
			ctx,
			generateValidatorLog("start token validation",
				"info",
				"RefreshTokenValidator",
			),
		)
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
				cacheManager.ClearUserRelatedCache(claims.Id.String())
				broker.GenerateLogEvent(
					ctx,
					generateValidatorLog(
						fmt.Sprintf("clearing cookies and session. Err - %s", err.Error()),
						"error",
						"RefreshTokenValidator",
					),
				)
				return err
			}

			// Retrieve cached refresh token
			cachedRefreshToken, err := cacheManager.RefreshToken(claims.Id.String())

			// Compare tokens
			if err != nil || cachedRefreshToken != refreshToken {
				generateValidatorLog(
					"Received token and token from cache don't match. Clearing cache & session",
					"error",
					"RefreshTokenValidator",
				)
				clearCookies(c)
				cacheManager.ClearUserRelatedCache(claims.Id.String())
				return goErrorHandler.NewError(
					goErrorHandler.ErrUnauthorized,
					fmt.Errorf("received token and token from cache don't match"),
				)
			}

			generateValidatorLog(
				"Refresh token has passed validation ",
				"info",
				"RefreshTokenValidator",
			)

			// Set user claims in the context for the next handler
			c.Set("user", claims)
			return next(c)
		}
	}
}
