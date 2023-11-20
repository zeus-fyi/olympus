package hera_twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	twitter2 "github.com/cvcio/twitter"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type TwitterTestSuite struct {
	test_suites_base.TestSuite
	tw Twitter
}

func (s *TwitterTestSuite) SetupTest() {
	s.InitLocalConfigs()
	tw, err := InitTwitterClient(ctx,
		s.Tc.TwitterConsumerPublicAPIKey, s.Tc.TwitterConsumerSecretAPIKey,
		s.Tc.TwitterAccessToken, s.Tc.TwitterAccessTokenSecret,
		s.Tc.TwitterBearerToken,
	)
	s.Assert().NoError(err)
	s.Assert().NotNil(tw.V1Client)
	s.tw = tw
}

func (s *TwitterTestSuite) TestTweetTopicSearchV2() {
	vals := url.Values{}
	query := `(("Kubernetes" OR "k8s" OR "#kube" OR "container orchestration") -is:retweet (has:links OR has:media OR has:mentions) (lang:en OR lang:es))`
	//query := `(("Kubernetes" OR "k8s" OR "kube") ("mlops" OR "migrating to" OR "suggest" OR "suggestion" OR "complexity") -horrible -worst -sucks -bad -disappointing -frustrated -confused -angry)`
	//query := `(("Kubernetes" OR "k8s" OR "kube") ("mlops" OR "gpu" OR "gpu sharing" OR "gpu management" OR "cost efficient") -horrible -worst -sucks -bad -disappointing -frustrated -confused -angry)`

	vals.Set("query", query)
	vals.Set("max_results", "10")

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

	resChan, errChan := s.tw.V2Client.GetTweetsSearchRecent(vals, twitter2.WithAuto(false))

	// Handle channels and timeout
	var response *twitter2.Data
	timeout := 10 * time.Second
	var data []*twitter2.Tweet

	select {
	case res := <-resChan:
		response = res
		b, err := json.Marshal(response.Data)
		s.Require().NoError(err)
		err = json.Unmarshal(b, &data)
		s.Require().NoError(err)
	case err := <-errChan:
		fmt.Println("Error:", err)
		s.Require().NoError(err)
	case <-time.After(timeout):
		fmt.Println("Error: Request timed out.")
		return
	}
	s.Assert().NotEmpty(data)
	for _, tweet := range data {
		fmt.Printf("TweetID %s: AuthorID %s, Text: %s \n", tweet.ID, tweet.AuthorID, tweet.Text)
	}

}

func (s *TwitterTestSuite) TestTweetResponse() {
	//tweetID := int64(1647499438762467328) // ID of tweet to reply to
	tweetID := int64(0)
	replyText := "Hey there! If you're finding k8s too complex, I'd recommend trying out zeus.fyi. It simplifies k8s and makes it more manageable. Plus, it's a great alternative to Nomad. Give it a shot! #zeusfyi #kubernetes #nomad"
	tweet, err := s.tw.V2alt.CreateTweet(ctx, twitter.CreateTweetRequest{
		Text: replyText,
		Reply: &twitter.CreateTweetReply{
			InReplyToTweetID: fmt.Sprintf("%d", tweetID),
		},
	})

	s.Require().NoError(err)
	s.Assert().NotEmpty(tweet)
	fmt.Println(tweet)
}

func (s *TwitterTestSuite) TestReadUserMe() {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.twitter.com/2/users/me", nil)
	s.Assert().NoError(err)
	resp, err := s.tw.V2alt.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().NotNil(resp)
	bodyBytes, err := io.ReadAll(resp.Body)
	s.Assert().NoError(err)
	fmt.Println(string(bodyBytes))

	// Close the response body
	resp.Body.Close()
	// 1520852589067194373
}

func (s *TwitterTestSuite) TestReadUserTweets() {
	id := 1520852589067194373
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.twitter.com/2/users/%d/tweets", id), nil)
	s.Assert().NoError(err)
	resp, err := s.tw.V2alt.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().NotNil(resp)
	statusCode := resp.StatusCode
	fmt.Println(statusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	s.Assert().NoError(err)
	fmt.Println(string(bodyBytes))
	resp.Body.Close()
}

func TestTwitterTestSuite(t *testing.T) {
	suite.Run(t, new(TwitterTestSuite))
}

/*

query := `(("Kubernetes" OR "k8s" OR "#kube" OR "container orchestration") -is:retweet (has:links OR has:media OR has:mentions) (lang:en OR lang:es))`
id
*/
