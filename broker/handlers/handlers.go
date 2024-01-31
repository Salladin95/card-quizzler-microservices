package handlers

import (
	"github.com/labstack/echo/v4"
)

const (
	SignInKey = "auth.sign-in.command"
	SignUpKey = "auth.sign-up.command"
)

func (bh *brokerHandlers) SignIn(c echo.Context) error {
	return bh.pushToQueueFromEndpoint(c, SignInKey)
}

func (bh *brokerHandlers) SignUp(c echo.Context) error {
	return bh.pushToQueueFromEndpoint(c, SignUpKey)
}
