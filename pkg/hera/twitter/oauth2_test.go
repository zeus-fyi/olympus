package hera_twitter

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/twitterv2"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	artemis_hydra_orchestrations_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

// // Set the redirect URI and scopes
// redirectURI := "https://hestia.zeus.fyi/oauth2/callback"
//
// // Create a code verifier and challenge
// codeVerifier, codeChallenge := createCodeVerifierAndChallenge()
//
// // OAuth2 config
//
//	conf := &oauth2.Config{
//		ClientID:     s.Tc.TwitterClientID,
//		ClientSecret: s.Tc.TwitterClientSecret, // Include this if your app is a confidential client
//		Scopes:       scopes,
//		RedirectURL:  redirectURI,
//		Endpoint: oauth2.Endpoint{
//			AuthURL:  "https://twitter.com/i/oauth2/authorize",
//			TokenURL: "https://api.twitter.com/2/oauth2/token",
//		},
//	}
//
// accessToken := addToken(s.Tc.TwitterConsumerPublicAPIKey, s.Tc.TwitterConsumerSecretAPIKey)
//
// // Resty client
// client := resty.New()
// client.SetAuthToken(accessToken)
// // Fetching user ID
// userID, err := fetchUserID(client)
//
//	if err != nil {
//		panic(err)
//	}
//
// // Bookmarking a tweet
// err = getBookmarks(client, userID, accessToken)
//
//	if err != nil {
//		panic(err)
//	}

func generateSecretKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *TwitterTestSuite) TestProvider() {
	goth.UseProviders(
		twitter.New(s.Tc.TwitterClientID, s.Tc.TwitterClientSecret, "http://localhost:9002/auth/twitter/callback"),
		twitterv2.New(s.Tc.TwitterClientID, s.Tc.TwitterClientSecret, "http://localhost:9002/auth/twitter/callback"),
	)

	ps := goth.GetProviders()
	s.Require().NotNil(ps)

	p, err := goth.GetProvider("twitterv2")
	s.Require().NoError(err)
	s.Require().NotNil(p)
	//cookie := &http.Cookie{
	//	Name:     aegis_sessions.SessionIDNickname,
	//	Value:    sessionID,
	//	HttpOnly: true,
	//	Secure:   true,
	//	Domain:   "zeus.fyi",
	//	SameSite: http.SameSiteNoneMode,
	//	Expires:  time.Now().Add(24 * time.Hour),
	//	Path:     "/",
	//}

	//s.Assert().NotNil(cookie)
	req := &http.Request{
		URL: &url.URL{},
	}
	q := req.URL.Query()
	q.Set("provider", "twitterv2")
	req.URL.RawQuery = q.Encode()

	pn, err := getProviderName(req)
	s.Require().NoError(err)
	s.Require().Equal("twitterv2", pn)

	req = httptest.NewRequest("GET", "http://localhost:9002/auth/twitter?provider=twitterv2", nil)

	sessionID := "191cc05e23d2b941ad16555fad5a403a2464987e134549f01b099b2d081c8e05"
	// Configure the cookie store for session management

	maxAge := 86400 * 30 // 30 days
	storeV := sessions.NewCookieStore([]byte(sessionID))
	storeV.MaxAge(maxAge)
	storeV.Options.Path = "/"
	storeV.Options.Domain = "zeus.fyi"
	storeV.Options.HttpOnly = true // HttpOnly should always be enabled
	storeV.Options.Secure = true

	gothic.Store = storeV

	// Prepare the request and response recorder
	req, _ = http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	// Set a value in the session
	key := "provider"
	value := "twitterv2" // Example provider name
	err = gothic.StoreInSession(key, value, req, res)
	s.NoError(err)

	// Now, simulate the client sending the response cookies as request cookies
	req.Header.Add("Cookie", res.Header().Get("Set-Cookie"))

	// Attempt to retrieve session value
	retrievedValue, err := gothic.GetFromSession(key, req)
	s.NoError(err)
	s.NotNil(retrievedValue)
	s.Equal(value, retrievedValue, "Retrieved session value should match the stored value")
}

func getProviderName(req *http.Request) (string, error) {

	// try to get it from the url param "provider"
	if p := req.URL.Query().Get("provider"); p != "" {
		return p, nil
	}

	// try to get it from the url param ":provider"
	if p := req.URL.Query().Get(":provider"); p != "" {
		return p, nil
	}

	// try to get it from the context's value of "provider" key
	if p, ok := mux.Vars(req)["provider"]; ok {
		return p, nil
	}

	//  try to get it from the go-context's value of "provider" key
	if p, ok := req.Context().Value("provider").(string); ok {
		return p, nil
	}

	return "", fmt.Errorf("could not find provider name")
}

func (s *TwitterTestSuite) TestOauth() {
	fmt.Println(generateSecretKey(32))

	awsAuthCfg := aegis_aws_auth.AuthAWS{
		AccountNumber: "",
		Region:        "us-west-1",
		AccessKey:     s.Tc.AwsAccessKeySecretManager,
		SecretKey:     s.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_auth.InitHydraSecretManagerAuthAWS(ctx, awsAuthCfg)
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(ctx, s.Ou, "twitter")
	s.Require().NoError(err)
	s.Require().NotNil(ps)

	tx, err := InitTwitterClient(ctx, ps.ConsumerPublic, ps.ConsumerSecret, ps.AccessTokenPublic, ps.AccessTokenSecret)
	s.Require().NoError(err)
	s.Require().NotNil(tx)

	veri, err := tx.V2Client.VerifyCredentials()
	s.Require().NoError(err)
	s.Require().NotNil(veri)

	urlWithParams := "https://api.twitter.com/2/users/me/bookmarks"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlWithParams, nil)
	s.Assert().NoError(err)
	resp, err := tx.V2Client.GetClient().Do(req)
	s.Assert().NoError(err)
	s.Assert().NotNil(resp)
	bodyBytes, err := io.ReadAll(resp.Body)
	s.Assert().NoError(err)
	fmt.Println(string(bodyBytes))

	//urlWithParams := "https://api.twitter.com/2/tweets/search/recent?" + vals.Encode()
	//
	//req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlWithParams, nil)
	//s.Assert().NoError(err)
	//resp, err := s.tw.V2alt.Client.Do(req)
	//s.Assert().NoError(err)
	//s.Assert().NotNil(resp)
	//bodyBytes, err := io.ReadAll(resp.Body)
	//s.Assert().NoError(err)
	//fmt.Println(string(bodyBytes))
	//
	//resp.Body.Close()

}
