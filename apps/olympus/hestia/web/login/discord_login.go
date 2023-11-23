package hestia_login

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

var (
	DiscordClientID     = ""
	DiscordRedirectURI  = ""
	DiscordClientSecret = ""
)

type DiscordLoginRequest struct {
	Credential string `json:"credential"`
}

func DiscordLoginHandler(c echo.Context) error {
	params := url.Values{}
	params.Add("client_id", DiscordClientID)
	params.Add("redirect_uri", DiscordRedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "identify guilds")

	authURL := "https://discord.com/api/oauth2/authorize?" + params.Encode()
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}
func (d *DiscordLoginRequest) VerifyDiscordLogin(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

// DiscordCallbackHandler handles the OAuth callback from Discord
func DiscordCallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, "Missing code")
	}
	// Prepare the request
	data := url.Values{}
	data.Set("client_id", DiscordClientID)
	data.Set("client_secret", DiscordClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", DiscordRedirectURI)
	// Exchange code for a token here
	// Implement the logic to get the token using the code
	req, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// Unmarshal the response
	var tokenResponse map[string]interface{}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return err
	}
	// Check if the response contains the access token
	if token, ok := tokenResponse["access_token"]; ok {
		return c.JSON(http.StatusOK, map[string]string{"token": token.(string)})
	}
	return c.JSON(http.StatusInternalServerError, "Failed to get access token")

}
