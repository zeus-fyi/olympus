package hestia_login

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hesta_base_test "github.com/zeus-fyi/olympus/hestia/api/test"
	"golang.org/x/oauth2"
)

type LoginTestSuite struct {
	hesta_base_test.HestiaBaseTestSuite
}

var ctx = context.Background()

func (t *LoginTestSuite) TestLogin() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)

	// https://hestia.zeus.fyi/auth/twitter/callback
	authorizeURL := "https://twitter.com/i/oauth2/authorize"
	tokenURL := "https://api.twitter.com/2/oauth2/token"
	conf := &oauth2.Config{
		RedirectURL:  "https://hestia.zeus.fyi/twitter/callback",
		ClientID:     t.Tc.TwitterClientID,
		ClientSecret: t.Tc.TwitterClientSecret,
		Scopes:       []string{"bookmark.write", "bookmark.read", "tweet.read", "users.read", "offline.access", "follows.read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizeURL,
			TokenURL: tokenURL,
		},
	}
	TwitterOAuthConfig = conf
	stateNonce := GenerateNonce()
	//verifier := GenerateCodeVerifier(128)
	challengeOpt := oauth2.SetAuthURLParam("code_challenge", PkCEChallengeWithSHA256(verifier))
	challengeMethodOpt := oauth2.SetAuthURLParam("code_challenge_method", "s256")
	redirectURL := TwitterOAuthConfig.AuthCodeURL(stateNonce, challengeOpt, challengeMethodOpt)
	fmt.Println(redirectURL)
}

func (t *LoginTestSuite) TestFetchToken() {
	// https://hestia.zeus.fyi/auth/twitter/callback
	authorizeURL := "https://twitter.com/i/oauth2/authorize"
	tokenURL := "https://api.twitter.com/2/oauth2/token"
	conf := &oauth2.Config{
		RedirectURL:  "https://hestia.zeus.fyi/twitter/callback",
		ClientID:     t.Tc.TwitterClientID,
		ClientSecret: t.Tc.TwitterClientSecret,
		Scopes:       []string{"bookmark.write", "bookmark.read", "tweet.read", "users.read", "offline.access", "follows.read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizeURL,
			TokenURL: tokenURL,
		},
	}
	TwitterOAuthConfig = conf

	to, err := FetchToken("bjZFdERlckJOX3I4WlMyWm1wZlZxU1RFc1JhWml1Q1VzQkMzd2RTSHVuckRoOjE3MDY3NjM2OTEwOTg6MTowOmFjOjE", "uvpjg4pN0ZphNPdsKUZahPCtfMnkGblvhT0M_Srp_-VZnTBAxXAPQ0ZynGRMAZfRTsIXpkTp1l9c7sI0c8FhA0q52efN0RM7DBUCMSSXz1y1-LkQ4LMi7f-jm_aFz7DzgtOfm1lcG_mKoXwPH4ss9j45RSfndb27JFjwDcy3Y2Y")
	t.Require().NoError(err)
	t.Require().NotNil(to)
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
