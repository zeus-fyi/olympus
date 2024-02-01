package hestia_login

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

var TwitterOAuthConfig = &oauth2.Config{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://twitter.com/i/oauth2/authorize",
		TokenURL: "https://api.twitter.com/2/oauth2/token",
	},
	RedirectURL: "http://localhost:9002/twitter/callback",
}

func CallbackHandler(c echo.Context) error {
	stateNonce := GenerateNonce()
	verifier := GenerateCodeVerifier(128)
	challengeOpt := oauth2.SetAuthURLParam("code_challenge", PkCEChallengeWithSHA256(verifier))
	challengeMethodOpt := oauth2.SetAuthURLParam("code_challenge_method", "s256")
	redirectURL := TwitterOAuthConfig.AuthCodeURL(stateNonce, challengeOpt, challengeMethodOpt)
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

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
	token, err := GenerateAccessToken(c.QueryParam("code"), c.QueryParam("verifier"))
	if err != nil {
		log.Err(err).Msg("TwitterCallbackHandler: Failed to generate access token")
		return err
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
