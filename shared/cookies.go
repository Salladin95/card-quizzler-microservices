package lib

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// ClearCookies removes refresh-token cookie from the provided echo.Context.
// It takes the echo.Context as parameter.
func ClearCookies(c echo.Context) {
	// Reset the specified cookies
	ResetCookies(c, "refresh-token")
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
