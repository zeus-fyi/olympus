package hestia_login

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func TwitterCallbackHandler(c echo.Context) error {
	// Complete the authentication process
	user, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		// Handle the error
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
