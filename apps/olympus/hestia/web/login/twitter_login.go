package hestia_login

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"github.com/rs/zerolog/log"
)

func CallbackHandler(c echo.Context) error {
	// Complete the authentication process
	providerName := c.Param("provider")
	if providerName == "twitter" {
		providerName = "twitterv2"
	}
	q := c.Request().URL.Query()
	q.Set("provider", providerName)
	c.Request().URL.RawQuery = q.Encode()
	user, err := CompleteUserAuth(c)
	if err != nil {
		log.Err(err).Interface("user", user).Msg("CallbackHandler: gothic.CompleteUserAuth")
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// Return user's data as JSON
	return c.JSON(http.StatusOK, user)
}

func LogoutHandler(c echo.Context) error {
	// Perform the logout using Goth
	gothic.Logout(c.Response().Writer, c.Request())
	// Redirect to the home page
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func AuthHandler(c echo.Context) error {
	log.Info().Msg("AuthHandler")
	// Try to complete the user authentication
	if gothUser, err := CompleteUserAuth(c); err == nil {
		// User is authenticated, return JSON
		return c.JSON(http.StatusOK, gothUser)
	} else {
		// Begin a new authentication process
		return BeginAuthHandler(c)
	}
}
