package hestia_login

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

var TwitterOAuthConfig = &oauth2.Config{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://twitter.com/i/oauth2/authorize",
		TokenURL: "https://api.twitter.com/2/oauth2/token",
	},
	Scopes:      []string{"bookmark.write", "bookmark.read", "tweet.read", "users.read", "offline.access", "follows.read"},
	RedirectURL: "https://hestia.zeus.fyi/twitter/callback",
}

// this is used to generate the URL to redirect the user to Twitter's OAuth2 login page

var verifier = GenerateCodeVerifier(128)

func CallbackHandler(c echo.Context) error {
	stateNonce := GenerateNonce()
	//verifier := GenerateCodeVerifier(128)
	challengeOpt := oauth2.SetAuthURLParam("code_challenge", PkCEChallengeWithSHA256(verifier))
	challengeMethodOpt := oauth2.SetAuthURLParam("code_challenge_method", "S256")
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
	//token, err := GenerateAccessToken(code, verifier)
	token, err := FetchToken(c.QueryParam("code"), verifier)
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

	// Prepare the URL values for the POST request body
	values := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {TwitterOAuthConfig.RedirectURL},
		"client_id":     {TwitterOAuthConfig.ClientID},
		"code_verifier": {codeVerifier},
	}

	// Log the values being sent in the request (ensure this does not log sensitive information in production)
	fmt.Printf("FetchToken: values=%v\n", values)

	// Create a Resty client
	client := resty.New()
	var token *oauth2.Token
	clientID := TwitterOAuthConfig.ClientID
	clientSecret := TwitterOAuthConfig.ClientSecret
	encodedCredentials := b64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Encode the client credentials using Base64
	authorizationHeaderValue := "Basic " + encodedCredentials

	// The "Authorization" header value should be "Basic " followed by the encoded credentials
	// Make the request using the Resty client
	resp, err := client.R().
		SetHeader("Authorization", authorizationHeaderValue). // Correctly set the Authorization header
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(values.Encode()).
		SetResult(&token).
		Post(TwitterOAuthConfig.Endpoint.TokenURL)

	fmt.Println(resp)
	if err != nil {
		return nil, err
	}

	return token, nil
}
