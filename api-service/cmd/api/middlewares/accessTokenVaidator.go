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

// AccessTokenValidator returns an Echo middleware that validates the access token.
func AccessTokenValidator(
	broker messageBroker.MessageBroker,
	cacheManager cacheManager.CacheManager,
	secret string,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			broker.GenerateLogEvent(
				ctx,
				generateValidatorLog("start token validation",
					"info",
					"AccessTokenValidator"),
			)

			// Get the Authorization header
			accessToken, err := ExtractAccessToken(c)
			if err != nil {
				return err
			}

			// Validate the access token
			claims, err := validateTokenString(accessToken, secret)
			if err != nil {
				broker.GenerateLogEvent(
					ctx,
					generateValidatorLog(
						err.Error(),
						"error",
						"AccessTokenValidator"),
				)
				return err
			}

			// Retrieve cached access token
			cachedAccessToken, err := cacheManager.AccessToken(claims.Id.String())
			// Compare tokens
			if err != nil || cachedAccessToken != accessToken {
				generateValidatorLog(
					"Received token and token from cache don't match. Clearing cache & session",
					"error",
					"AccessTokenValidator")
				clearCookies(c)
				cacheManager.ClearDataByUID(claims.Id.String())
				return goErrorHandler.NewError(
					goErrorHandler.ErrUnauthorized,
					fmt.Errorf("received access token and cached token don't match"),
				)
			}

			generateValidatorLog(
				"Access token has passed validation",
				"info",
				"AccessTokenValidator")
			c.Set("user", claims)
			return next(c)
		}
	}
}
