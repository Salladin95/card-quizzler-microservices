package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func RefreshTokenValidator(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		log.Info("RefreshTokenValidator - start token validation")
		return func(c echo.Context) error {
			tokenString, err := ExtractRefreshToken(c)
			if err != nil {
				return err
			}
			claims, err := validateTokenString(tokenString, secret)
			if err != nil {
				clearSessionAndCookies(c)
				return err
			}
			// TODO: EXTRACT TOKEN FROM SESSION AND COMPARE THEM
			//ctx := c.Request().Context()
			log.Errorf("****** refersh token has passed validation *******\n")

			c.Set("user", claims)
			return next(c)
		}
	}
}
