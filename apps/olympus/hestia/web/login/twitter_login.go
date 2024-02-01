package hestia_login

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

var TwitterOAuthConfig = &oauth2.Config{}

func CallbackHandler(c echo.Context) error {
	redirectURL := TwitterOAuthConfig.AuthCodeURL(state)
	log.Info().Msgf("Redirecting to %s", redirectURL)
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func LogoutHandler(c echo.Context) error {
	return nil
}

func AuthHandler(c echo.Context) error {
	return nil
}
