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
	if gothUser, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request()); err == nil {
		// User is authenticated, return JSON
		log.Info().Interface("gothUser", gothUser).Msg("CallbackHandler")
		return c.JSON(http.StatusOK, gothUser)
	}
	url, err := gothic.GetAuthURL(c.Response().Writer, c.Request())
	if err != nil {
		log.Err(err).Interface("url", url).Msg("CallbackHandlerError")
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	log.Info().Interface("url", url).Msg("CallbackHandlerGetAuthURL")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func LogoutHandler(c echo.Context) error {
	// Perform the logout using Goth
	gothic.Logout(c.Response().Writer, c.Request())
	// Redirect to the home page
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func AuthHandler(c echo.Context) error {
	log.Info().Msg("AuthHandler")
	providerName := c.Param("provider")
	if providerName == "twitter" {
		providerName = "twitterv2"
	}
	q := c.Request().URL.Query()
	q.Set("provider", providerName)
	c.Request().URL.RawQuery = q.Encode()
	if gothUser, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request()); err == nil {
		// User is authenticated, return JSON
		log.Info().Interface("gothUser", gothUser).Msg("CallbackHandler")
		return c.JSON(http.StatusOK, gothUser)
	}
	url, err := gothic.GetAuthURL(c.Response().Writer, c.Request())
	if err != nil {
		log.Err(err).Interface("url", url).Msg("CallbackHandler")
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	log.Info().Interface("url", url).Msg("CallbackHandler")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}
