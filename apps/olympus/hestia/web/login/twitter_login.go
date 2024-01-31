package hestia_login

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"github.com/rs/zerolog/log"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

var Providers = ProviderIndex{
	Providers:    []string{},
	ProvidersMap: make(map[string]string),
}

func TwitterCallbackHandler(c echo.Context) error {
	// Complete the authentication process
	user, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		log.Err(err).Msg("TwitterCallbackHandler: gothic.CompleteUserAuth")
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// Return user's data as JSON
	return c.JSON(http.StatusOK, user)
}

func TwitterLogoutHandler(c echo.Context) error {
	// Perform the logout using Goth
	gothic.Logout(c.Response().Writer, c.Request())
	// Redirect to the home page
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func TwitterAuthHandler(c echo.Context) error {
	// Try to complete the user authentication
	if gothUser, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request()); err == nil {
		// User is authenticated, return JSON
		return c.JSON(http.StatusOK, gothUser)
	} else {
		// Begin a new authentication process
		gothic.BeginAuthHandler(c.Response().Writer, c.Request())
		return nil
	}
}
