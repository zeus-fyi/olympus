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

	to, err := FetchToken("N1ctek9XbGpkZWs2aW5vSmV5Z2ZxdzR0Qm4wcEVleDQwVWhwWUozT2RMTWp4OjE3MDY3NjM5NTc2ODk6MTowOmFjOjE", "_S_An_a_VEh3atLl7hfAaiYtTm4IZIXqF_aY5P3yJi8tICGmNaO2mXei-uqqhdWvyTCm2PeBp8OoiEWQq7jCELoSsefhnU0c4fKh_3tu_ZFqOSan9FTxN9Qc_LXW0H3gs3kp9KgFG2J8ZrICNfzqVCoXJ4OM0_rVOt14XxcCA7Y")
	t.Require().NoError(err)
	t.Require().NotNil(to)
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
