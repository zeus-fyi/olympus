package hestia_login

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
)

func RedditLoginHandler(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, Conf.AuthCodeURL(state))
}

func RedditCallbackHandler(c echo.Context) error {
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
		ts, aok := token.(string)
		if !aok {
			return c.JSON(http.StatusInternalServerError, "Failed to get access token")
		}
		nk := create_keys.NewCreateKey(internalOrgID, ts)
		nk.PublicKeyVerified = true
		nk.PublicKeyName = "discord"
		nk.CreatedAt = time.Now()
		nk.UserID = 7138958574876245567
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, "ok")
	}

	return c.JSON(http.StatusInternalServerError, nil)

}
