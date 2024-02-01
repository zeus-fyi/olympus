package hestia_login

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"github.com/rs/zerolog/log"
)

func CallbackHandler(c echo.Context) error {
	providerName := c.Param("provider")
	if providerName == "twitter" {
		providerName = "twitterv2"
	}
	q := c.Request().URL.Query()
	q.Set("provider", providerName)
	c.Request().URL.RawQuery = q.Encode()

	// Attempt to complete the user authentication process
	if gothUser, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request()); err == nil {
		// User is authenticated, return JSON
		log.Info().Interface("gothUser", gothUser).Msg("CallbackHandler: User authenticated")
		return c.JSON(http.StatusOK, gothUser)
	} else {
		// If authentication fails, log the error and redirect with an error message
		log.Error().Err(err).Msg("CallbackHandler: Authentication failed")
		url, authErr := gothic.GetAuthURL(c.Response().Writer, c.Request())
		if authErr != nil {
			errorMessage := fmt.Sprintf("Authentication failed: %s", authErr.Error())
			log.Err(authErr).Interface("url", url).Msg("CallbackHandlerError: Failed to get auth URL")
			// Consider redirecting to an error page or passing the error message to the frontend
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": errorMessage})
		}
		log.Info().Interface("url", url).Msg("CallbackHandler: Redirecting to auth URL")
		// Optionally, append a query parameter to the URL with an error message
		return c.Redirect(http.StatusTemporaryRedirect, url+"?error=authentication_failed")
	}
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
	gothUser, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		log.Error().Err(err).Msg("AuthHandler: Error completing user authentication")
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, gothUser)
}
