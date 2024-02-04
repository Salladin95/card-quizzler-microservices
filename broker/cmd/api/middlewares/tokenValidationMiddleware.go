package middlewares

import (
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"strings"
)

func TokenValidationMiddleWare(authClient *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("****** START VALIDATING TOKEN ***********")
			ctx := c.Request().Context()

			// Get the Authorization header
			authHeader := c.Request().Header.Get("Authorization")

			// Check if the header is present and starts with "Bearer "
			if authHeader == "" {
				return goErrorHandler.Unauthorized()
			}

			bearerToken := strings.Split(authHeader, "Bearer ")[1]
			fmt.Println(bearerToken)
			fmt.Println()

			token, err := authClient.VerifyIDToken(ctx, bearerToken)
			if err != nil {
				fmt.Printf("****** failed token validation: %v ********\n", err)
				return goErrorHandler.Unauthorized()
			}

			fmt.Printf("verified token - %v\n", token)

			fmt.Println("!!!!!!!!!!!!!!!!!! TOKEN VALIDATION IS PASSED !!!!!!!!!!!!!!!!!!")
			return next(c)

		}
	}
}
