package middlewares

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/cacheManager"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// AccessTokenValidator returns an Echo middleware that validates the access token.
func AccessTokenValidator(cacheManager cacheManager.CacheManager, secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("****** START VALIDATING ACCESS TOKEN ***********")

			// Get the Authorization header
			accessToken, err := ExtractAccessToken(c)
			if err != nil {
				return err
			}

			// Validate the access token
			claims, err := validateTokenString(accessToken, secret)
			if err != nil {
				return err
			}

			// Retrieve cached access token
			cachedAccessToken, err := cacheManager.AccessToken(claims.Id.String())
			// Compare tokens
			if err != nil || cachedAccessToken != accessToken {
				log.Infof("********* Received access token and cached token don't match ***********")
				clearCookies(c)
				cacheManager.ClearUserData(claims.Id.String())
				return goErrorHandler.NewError(
					goErrorHandler.ErrUnauthorized,
					fmt.Errorf("received access token and cached token don't match"),
				)
			}

			log.Errorf("****** Access token has passed validation *******\n")
			c.Set("user", claims)
			return next(c)
		}
	}
}
