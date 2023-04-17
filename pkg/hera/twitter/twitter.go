package hera_twitter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cvcio/twitter"
	twitter_auth "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	v2twitter "github.com/g8rswimmer/go-twitter/v2"
	"github.com/rs/zerolog/log"
)

type Twitter struct {
	V2Client *twitter.Twitter
	V2alt    *v2twitter.Client
	V1Client *twitter_auth.Client
}

const (
	BaseURL            = "https://api.twitter.com/2"
	RequestTokenURL    = "https://api.twitter.com/oauth/request_token"
	AuthorizeTokenURL  = "https://api.twitter.com/oauth/authorize"
	AccessTokenURL     = "https://api.twitter.com/oauth/access_token"
	TokenURL           = "https://api.twitter.com/oauth2/token"
	RateLimitStatusURL = "https://api.twitter.com/1.1/application/rate_limit_status.json"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func InitTwitterClient(ctx context.Context, consumerKey, consumerSecret, accessToken, accessSecret, bearer string) (Twitter, error) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(ctx, token)
	client := twitter_auth.NewClient(httpClient)
	api, err := twitter.NewTwitter(consumerKey, consumerSecret)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error creating twitter client")
		return Twitter{}, err
	}

	// Create HTTP client
	httpClient = config.Client(oauth1.NoContext, token)
	// Create Twitter client using OAuth1 authentication
	// http.Client will automatically authorize Requests
	v2C := &v2twitter.Client{
		Authorizer: authorize{
			Token: token.Token,
		},
		Client: httpClient,
		Host:   "https://api.twitter.com",
	}
	return Twitter{V1Client: client, V2Client: api, V2alt: v2C}, err
}
