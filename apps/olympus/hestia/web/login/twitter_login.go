package hestia_login

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var TwitterOAuthConfig = &oauth2.Config{}

func CallbackHandler(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, TwitterOAuthConfig.AuthCodeURL(state))
}

func LogoutHandler(c echo.Context) error {
	return nil
}

func AuthHandler(c echo.Context) error {
	return nil
}
