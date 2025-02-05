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
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/iris/orchestrations"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
	"golang.org/x/oauth2"
)

var TwitterOAuthConfig = &oauth2.Config{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://twitter.com/i/oauth2/authorize",
		TokenURL: "https://api.twitter.com/2/oauth2/token",
	},
	Scopes: []string{
		"bookmark.write",
		"bookmark.read",
		"tweet.read",
		"tweet.write",
		"tweet.moderate.write",
		"users.read",
		"follows.read",
		"follows.write",
		"offline.access",
		"space.read",
		"mute.read",
		"mute.write",
		"like.read",
		"like.write",
		"list.read",
		"list.write",
		"block.read",
		"block.write",
	},
	RedirectURL: "https://cloud.zeus.fyi/social/v1/twitter/callback",
}

const (
	RedirectBackPlatform = "https://cloud.zeus.fyi/social/v1/twitter/callback"
)

//	[]string{"bookmark.write", "bookmark.read", "tweet.read", "users.read", "offline.access", "follows.read"},
//
// this is used to generate the URL to redirect the user to Twitter's OAuth2 login page
var (
	ch    = cache.New(5*time.Minute, 10*time.Minute)
	chOrg = cache.New(5*time.Minute, 10*time.Minute)
)

func CallbackHandler(c echo.Context) error {
	provider := c.Param("provider")
	switch provider {
	case "reddit":
		return RedditCallbackHandler(c)
	case "twitter":
	}
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	stateNonce := GenerateNonce()
	verifier := GenerateCodeVerifier(128)
	codeChallenge := PkCEChallengeWithSHA256(verifier)
	//log.Info().Str("codeChallenge", codeChallenge).Interface("stateNonce", stateNonce).Interface("verifier", verifier).Msg("START CallbackHandler: Callback")
	// Store the verifier using stateNonce as the key
	//log.Info().Interface("orgUser", ou).Str("stateNonce", stateNonce).Msg("BEGIN: Storing stateNonce in cache")
	ch.Set(stateNonce, verifier, cache.DefaultExpiration)
	chOrg.Set(stateNonce, fmt.Sprintf("%d", ou.OrgID), cache.DefaultExpiration)

	challengeOpt := oauth2.SetAuthURLParam("code_challenge", codeChallenge)
	challengeMethodOpt := oauth2.SetAuthURLParam("code_challenge_method", "S256")
	redirectURL := TwitterOAuthConfig.AuthCodeURL(stateNonce, challengeOpt, challengeMethodOpt)
	return c.JSON(http.StatusOK, redirectURL)
}

func TwitterCallbackHandler(c echo.Context) error {
	//log.Printf("START TwitterCallbackHandler Callback: Method=%s, URL=%s", c.Request().Method, c.Request().URL)

	code := c.QueryParam("code")
	stateNonce := c.QueryParam("state")
	//log.Info().Interface("stateNonce", stateNonce).Msg("TwitterCallbackHandler: Handling Twitter Callback")
	//log.Info().Str("code", code).Str("state", stateNonce).Msg("TwitterCallbackHandler: Handling Twitter Callback")
	verifier, found := ch.Get(stateNonce)
	if !found {
		log.Warn().Msg("TwitterCallbackHandler: Failed to retrieve verifier from cache")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	verifierStr, ok := verifier.(string)
	if !ok {
		log.Err(fmt.Errorf("TwitterCallbackHandler: Failed to cast verifier to string")).Msg("TwitterCallbackHandler: Failed to cast verifier to string")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	orgStrID, found := chOrg.Get(stateNonce)
	if !found {
		log.Warn().Msg("TwitterCallbackHandler: Failed to retrieve verifier from cache")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	orgStrIDFromCache, ok := orgStrID.(string)
	if !ok {
		log.Info().Interface("orgStrID", orgStrID).Msg("TwitterCallbackHandler: orgID from cache")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	log.Info().Interface("orgStrIDFromCache", orgStrIDFromCache).Msg("TwitterCallbackHandler: orgID from cache")
	log.Info().Interface("chOrg", chOrg).Msg("TwitterCallbackHandler: org cache")
	ouid, err := strconv.Atoi(orgStrIDFromCache)
	if err != nil {
		log.Err(err).Msg("TwitterCallbackHandler: Failed to cast orgID to int")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ou := org_users.OrgUser{}
	ou.OrgID = ouid
	c.Set("orgUser", ou)
	log.Info().Interface("orgUser", ou).Msg("TwitterCallbackHandler: Successfully retrieved orgID from cache")
	token, err := FetchToken(code, verifierStr)
	if err != nil {
		log.Err(err).Msg("TwitterCallbackHandler: Failed to generate access token")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	log.Info().Interface("token", token).Msg("TwitterCallbackHandler: Successfully generated access token")
	if token == nil {
		log.Warn().Msg("TwitterCallbackHandler: Token is nil")
		//return c.JSON(http.StatusInternalServerError, "Token is nil")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	//log.Info().Interface("token", token).Msg("TwitterCallbackHandler: Successfully generated access token")
	//return c.Redirect(http.StatusTemporaryRedirect, "https://cloud.zeus.fyi/ai")
	r := resty_base.GetBaseRestyClient("https://api.twitter.com", token.AccessToken)
	tm := TwitterMe{}
	resp, err := r.R().SetResult(&tm).Get("/2/users/me")
	if err != nil {
		log.Err(err).Interface("resp", resp).Interface("tm", tm).Msg("TwitterCallbackHandler: Failed to fetch user data")
		return c.JSON(http.StatusInternalServerError, "Failed to fetch user data")
	}
	log.Info().Interface("tm", tm).Msg("TwitterCallbackHandler: TwitterMe")
	// username is the unique handle identifier for the user
	sr := SecretsRequest{
		Name:  fmt.Sprintf("api-twitter-%s", tm.Data.Username),
		Key:   "mockingbird",
		Value: token.AccessToken,
	}
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:           ou,
		OrgGroupName: fmt.Sprintf("twitter-%s", tm.Data.Username),
		Routes:       []string{"https://api.twitter.com/2/"},
	}
	ctx := context.Background()
	err = platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisPlatformSetupRequestWorkflow(ctx, ipr)
	if err != nil {
		log.Err(err).Msg("TwitterCallbackHandler: CreateOrgGroupRoutesRequest")
		return err
	}
	return sr.CreateOrUpdateSecret(c, false)
}

type TwitterMe struct {
	Data struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"data"`
}

func FetchToken(code string, codeVerifier string) (*oauth2.Token, error) {
	// Prepare the URL values for the POST request body
	values := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {RedirectBackPlatform},
		"client_id":     {TwitterOAuthConfig.ClientID},
		"code_verifier": {codeVerifier},
	}

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

	if err != nil {
		log.Err(err).Interface("statusCode", resp.StatusCode()).Msg("FetchToken: Failed to fetch token")
		return nil, err
	}
	//log.Info().Interface("resp", resp.Body()).Msg("FetchToken: Successfully fetched token")
	return token, nil
}

func LogoutHandler(c echo.Context) error {
	return nil
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
