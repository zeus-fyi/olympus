package hera_twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/cvcio/twitter"
	twitter2 "github.com/cvcio/twitter"
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

var TwitterClient Twitter

func InitTwitterClient(ctx context.Context, consumerKey, consumerSecret, accessToken, accessSecret string) (Twitter, error) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(ctx, token)
	client := twitter_auth.NewClient(httpClient)
	api, err := twitter.NewTwitter(consumerKey, consumerSecret)
	if err != nil {
		log.Err(err).Msg("error creating twitter client")
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
	TwitterClient = Twitter{V1Client: client, V2Client: api, V2alt: v2C}
	return Twitter{V1Client: client, V2Client: api, V2alt: v2C}, err
}

func InitOrgTwitterClient(ctx context.Context, consumerKey, consumerSecret string) (Twitter, error) {
	api, err := twitter.NewTwitter(consumerKey, consumerSecret)
	if err != nil {
		log.Err(err).Msg("error creating twitter client")
		return Twitter{}, err
	}
	return Twitter{V2Client: api}, err
}

func (tc *Twitter) GetTweets(ctx context.Context, query string, maxResults, maxTweetID int) ([]*twitter2.Tweet, error) {
	vals := url.Values{}
	vals.Set("query", query)
	vals.Set("max_results", fmt.Sprintf("%d", maxResults))
	vals.Set("since_id", fmt.Sprintf("%d", maxTweetID))

	resChan, errChan := tc.V2Client.GetTweetsSearchRecent(vals, twitter2.WithAuto(false))

	// Handle channels and timeout
	var response *twitter2.Data
	timeout := 10 * time.Second
	var data []*twitter2.Tweet

	select {
	case res := <-resChan:
		response = res
		b, err := json.Marshal(response.Data)
		if err != nil {
			log.Err(err).Msg("error marshalling response data")
			return nil, err
		}
		err = json.Unmarshal(b, &data)
		if err != nil {
			log.Err(err).Msg("error marshalling response data")
			return nil, err
		}
	case err := <-errChan:
		if err != nil {
			log.Err(err).Msg("error marshalling response data")
			return nil, err
		}
	case <-time.After(timeout):
		err := fmt.Errorf("error: Request timed out.")
		if err != nil {
			log.Err(err).Msg("error marshalling response data")
			return nil, err
		}
	}
	return data, nil
}

func GetTweets(ctx context.Context, query string, maxResults, maxTweetID int) ([]*twitter2.Tweet, error) {
	return nil, nil
}
