package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func AccessTokenValidator(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("****** START VALIDATING TOKEN ***********")
			// Get the Authorization header
			accessToken, err := ExtractAccessToken(c)
			if err != nil {
				return err
			}
			claims, err := validateTokenString(accessToken, secret)
			if err != nil {
				return err
			}
			// TODO: EXTRACT TOKEN FROM SESSION AND COMPARE THEM
			//ctx := c.Request().Context()
			log.Errorf("****** access token has passed validation *******\n")
			c.Set("user", claims)
			return next(c)
		}
	}
}
