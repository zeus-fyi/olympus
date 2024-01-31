package hera_twitter

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"

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
