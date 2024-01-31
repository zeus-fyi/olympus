package hestia_login

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

//func TwitterLoginHandler(c echo.Context) error {
//	return c.Redirect(http.StatusTemporaryRedirect, Conf.AuthCodeURL(state))
//}

//ps, err := aws_secrets.GetMockingbirdPlatformSecrets(c.Request().Context(), s.Ou, "twitter")

func TwitterCallbackHandler(c echo.Context) error {

	// Complete the authentication process
	user, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		// Handle the error, perhaps by redirecting to an error page
		return c.String(http.StatusInternalServerError, err.Error())
	}
	// At this point, authentication is successful, and you have the user data
	// You can now perform actions like creating a user session, storing user data, etc.
	// For demonstration, let's just return the user's data as JSON
	return c.JSON(http.StatusOK, user)
}

func TwitterLogoutHandler(c echo.Context) error {
	req := c.Request()
	res := c.Response().Writer

	// Perform the logout using Goth
	gothic.Logout(res, req)

	// Set the header and redirect
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func TwitterAuthHandler(c echo.Context) error {
	req := c.Request()
	res := c.Response().Writer

	// Try to complete the user authentication
	gothUser, err := gothic.CompleteUserAuth(res, req)
	if err == nil {
		// User is authenticated, render the success template
		return c.JSON(http.StatusOK, gothUser)
	} else {
		// Begin a new authentication process
		gothic.BeginAuthHandler(res, req)
		return nil
	}
}
