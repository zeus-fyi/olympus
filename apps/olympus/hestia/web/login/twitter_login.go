package hestia_login

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
)

//func TwitterLoginHandler(c echo.Context) error {
//	return c.Redirect(http.StatusTemporaryRedirect, Conf.AuthCodeURL(state))
//}

func TwitterCallbackHandler(c echo.Context) error {

	//ps, err := aws_secrets.GetMockingbirdPlatformSecrets(c.Request().Context(), s.Ou, "twitter")

	goth.UseProviders(
		twitter.New("TWITTER_KEY", "TWITTER_SECRET", "http://localhost:9002/auth/twitter/callback"),
	)
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

func TwitterAuthHandler(c echo.Context) error {
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
