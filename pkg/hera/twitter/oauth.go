package hera_twitter

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

func randomBytesInHex(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("could not generate %d random bytes: %v", count, err)
	}

	return hex.EncodeToString(buf), nil
}

type AuthURL struct {
	URL          string
	State        string
	CodeVerifier string
}

func (u *AuthURL) String() string {
	return u.URL
}

func AuthorizeURL(config *oauth2.Config) (*AuthURL, error) {
	codeVerifier, verifierErr := randomBytesInHex(32) // 64 character string here
	if verifierErr != nil {
		return nil, fmt.Errorf("could not create a code verifier: %v", verifierErr)
	}
	sha2 := sha256.New()
	io.WriteString(sha2, codeVerifier)
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))

	stateRand, stateErr := randomBytesInHex(24)
	if stateErr != nil {
		return nil, fmt.Errorf("could not generate random state: %v", stateErr)
	}

	authUrl := config.AuthCodeURL(
		stateRand,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	return &AuthURL{
		URL:          authUrl,
		State:        state,
		CodeVerifier: codeVerifier,
	}, nil
}

const (
	QUERY_STATE = "state"
	QUERY_CODE  = "code"
)

type OAuthRedirectHandler struct {
	State        string
	CodeVerifier string
	OAuthConfig  *oauth2.Config
}

func textResponse(rw http.ResponseWriter, status int, body string) {
	rw.Header().Add("Content-Type", "text/plain")
	rw.WriteHeader(status)
	io.WriteString(rw, body)
}

func (h *OAuthRedirectHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	stateVal := query.Get(QUERY_STATE)
	// prevent timing attacks on state
	if subtle.ConstantTimeCompare([]byte(h.State), []byte(stateVal)) == 0 {
		textResponse(rw, http.StatusBadRequest, "Invalid State")
		return
	}

	code := query.Get(QUERY_CODE)
	if code == "" {
		textResponse(rw, http.StatusBadRequest, "Missing Code")
		return
	}

	token, err := h.OAuthConfig.Exchange(
		request.Context(),
		code,
		oauth2.SetAuthURLParam("code_verifier", h.CodeVerifier),
	)
	if err != nil {
		textResponse(rw, http.StatusInternalServerError, err.Error())
		return
	}

	// probably do something more legit with this token...
	textResponse(rw, http.StatusOK, token.AccessToken)
}
