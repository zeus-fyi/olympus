package hestia_login

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

var TwitterOAuthConfig = &oauth2.Config{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://twitter.com/i/oauth2/authorize",
		TokenURL: "https://api.twitter.com/2/oauth2/token",
	},
	RedirectURL: "https://hestia.zeus.fyi/twitter/callback",
}

// this is used to generate the URL to redirect the user to Twitter's OAuth2 login page

var verifier = GenerateCodeVerifier(128)

func CallbackHandler(c echo.Context) error {
	stateNonce := GenerateNonce()
	//verifier := GenerateCodeVerifier(128)
	challengeOpt := oauth2.SetAuthURLParam("code_challenge", PkCEChallengeWithSHA256(verifier))
	challengeMethodOpt := oauth2.SetAuthURLParam("code_challenge_method", "s256")
	redirectURL := TwitterOAuthConfig.AuthCodeURL(stateNonce, challengeOpt, challengeMethodOpt)
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// this is used to generate the access token from the code returned by Twitter's OAuth2 login page

func GenerateAccessToken(code string, verifier string) (*oauth2.Token, error) {
	ctx := context.Background()
	verifierOpt := oauth2.SetAuthURLParam("code_verifier", verifier)
	token, err := TwitterOAuthConfig.Exchange(ctx, code, verifierOpt)
	if err != nil {
		log.Err(err).Msg("GenerateAccessToken: Failed to exchange code for access token")
		return nil, err
	}
	return token, nil
}

func TwitterCallbackHandler(c echo.Context) error {
	log.Printf("Handling Twitter Callback: Method=%s, URL=%s", c.Request().Method, c.Request().URL)

	code := c.QueryParam("code")
	log.Info().Msgf("TwitterCallbackHandler: code=%s", code)
	token, err := FetchToken(c.QueryParam("code"), verifier)
	//token, err := GenerateAccessToken(c.QueryParam("code"), verifier)
	if err != nil {
		log.Err(err).Msg("TwitterCallbackHandler: Failed to generate access token")
		return c.JSON(http.StatusInternalServerError, "Failed to generate access token")
	}
	return c.JSON(http.StatusOK, token)
}

// PkCEChallengeWithSHA256 base64-URL-encoded SHA256 hash of verifier, per rfc 7636
func PkCEChallengeWithSHA256(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	challenge := b64.RawURLEncoding.EncodeToString(sum[:])
	return challenge
}

// GenerateCodeVerifier Generate code verifier (length 43~128) for PKCE.
func GenerateCodeVerifier(length int) string {
	if length > 128 {
		length = 128
	}
	if length < 43 {
		length = 43
	}
	return randStringBytes(length)
}

// GenerateNonce Generate random nonce.
func GenerateNonce() string {
	return b64.RawURLEncoding.EncodeToString([]byte(randStringBytes(8)))
}

// randStringBytes Return random string by length
func randStringBytes(n int) string {
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		return ""
	}
	return b64.RawURLEncoding.EncodeToString(b[:])
}

func LogoutHandler(c echo.Context) error {
	return nil
}
func FetchToken(code string, codeVerifier string) (*oauth2.Token, error) {
	ctx := context.Background()

	// Prepare the URL values for the POST request body
	values := url.Values{
		"code":          []string{code},
		"grant_type":    []string{"authorization_code"},
		"redirect_uri":  []string{TwitterOAuthConfig.RedirectURL},
		"client_id":     []string{TwitterOAuthConfig.ClientID},
		"code_verifier": []string{codeVerifier},
	}

	fmt.Println(values)
	// If the client_secret is not empty, include it in the request
	if TwitterOAuthConfig.ClientSecret != "" {
		values.Set("client_secret", TwitterOAuthConfig.ClientSecret)
	}

	log.Info().Msgf("FetchToken: values=%v", values)
	// Create a request to send to the token endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, TwitterOAuthConfig.Endpoint.TokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		log.Err(err).Msg("oauth2: cannot create request")
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	log.Printf("Handling Twitter FetchToken: Method=%s, URL=%s", req.Method, req.URL)

	// Make the request using the http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success
	if resp.StatusCode != http.StatusOK {
		log.Warn().Interface("resp.Body", resp.Body).Msg("oauth2: server response indicates failure")
		return nil, errors.New("oauth2: server response indicates failure")
	}

	// Parse the response body to extract the token
	var token *oauth2.Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
