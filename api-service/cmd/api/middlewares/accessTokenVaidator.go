package middlewares

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
)

// AccessTokenValidator returns an Echo middleware that validates the access token.
func AccessTokenValidator(
	cacheManager cacheManager.CacheManager,
	secret string,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			lib.LogInfo("Start access token validation")

			// Get the Authorization header
			accessToken, err := lib.ExtractAccessToken(c)
			if err != nil {
				return err
			}

			// Validate the access token
			claims, err := lib.ValidateTokenString(accessToken, secret)
			if err != nil {
				lib.LogError(
					err,
				)
				return err
			}

			// Retrieve cached access token
			cachedAccessToken, err := cacheManager.AccessToken(claims.Id)
			// Compare tokens
			if err != nil || cachedAccessToken != accessToken {

				lib.LogInfo("Received token and token from cache don't match. Clearing cache & session")

				lib.ClearCookies(c)
				cacheManager.ClearUserRelatedCache(claims.Id)
				return goErrorHandler.NewError(
					goErrorHandler.ErrUnauthorized,
					fmt.Errorf("received access token and cached token don't match"),
				)
			}

			lib.LogInfo("Access token has passed validation")

			c.Set("user", claims)
			return next(c)
		}
	}
}
