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
	"github.com/dghubble/go-twitter/twitter"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

var ctx = context.Background()

type TwitterTestSuite struct {
	test_suites_base.TestSuite
	tw Twitter
}

func (s *TwitterTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	tw, err := InitPkgTwitterClient(ctx,
		s.Tc.TwitterConsumerPublicAPIKey, s.Tc.TwitterConsumerSecretAPIKey,
		s.Tc.TwitterAccessToken, s.Tc.TwitterAccessTokenSecret,
	)
	s.Assert().NoError(err)
	s.Assert().NotNil(tw.V1Client)
	s.tw = tw

}

func (s *TwitterTestSuite) TestOauth1() {
	tc, terr := InitOrgTwitterClient(ctx, s.Tc.TwitterConsumerPublicAPIKey, s.Tc.TwitterConsumerSecretAPIKey)
	s.Require().NoError(terr)
	s.Assert().NotNil(tc.V2Client)

	res, resp, err := tc.V1Client.Users.Lookup(&twitter.UserLookupParams{
		ScreenName: []string{"ctrl_alt_lulz"},
	})
	s.Require().NoError(err)
	s.Assert().NotNil(resp)
	s.Assert().NotNil(res)

}

func (s *TwitterTestSuite) TestTweetTopicSearchOrgV2() {
	tc, terr := InitOrgTwitterClient(ctx, s.Tc.TwitterConsumerPublicAPIKey, s.Tc.TwitterConsumerSecretAPIKey)
	s.Require().NoError(terr)
	s.Assert().NotNil(tc.V2Client)
	id := "1520852589067194373"
	vals := url.Values{}
	vals.Set("max_results", "10")

	resChan, errChan := tc.V2Client.GetUserTweets(id, vals, twitter2.WithAuto(false))

	// Close the response body

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

func (s *TwitterTestSuite) TestTweetTopicSearchV2() {
	vals := url.Values{}
	query := `(("Kubernetes" OR "k8s" OR "#kube" OR "container orchestration") -is:retweet (has:links OR has:media OR has:mentions) (lang:en))`
	//query := `(("Kubernetes" OR "k8s" OR "kube") ("mlops" OR "migrating to" OR "suggest" OR "suggestion" OR "complexity") -horrible -worst -sucks -bad -disappointing -frustrated -confused -angry)`
	//query := `(("Kubernetes" OR "k8s" OR "kube") ("mlops" OR "gpu" OR "gpu sharing" OR "gpu management" OR "cost efficient") -horrible -worst -sucks -bad -disappointing -frustrated -confused -angry)`

	vals.Set("query", query)
	vals.Set("max_results", "10")
	vals.Set("since_id", "1726707774195728627")

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

	resp, err := hera_search.InsertIncomingTweets(ctx, 1700514815519783000, data)
	s.Require().NoError(err)
	s.Assert().NotEmpty(resp)
}

func (s *TwitterTestSuite) TestApiGets() {
	vals := url.Values{}
	vals.Set("screen_name", "OpenAI")

	urlWithParams := "https://api.twitter.com/2/users/me/bookmarks"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlWithParams, nil)
	s.Assert().NoError(err)
	resp, err := s.tw.V2alt.Client.Do(req)
	s.Assert().NoError(err)
	s.Assert().NotNil(resp)
	bodyBytes, err := io.ReadAll(resp.Body)
	s.Assert().NoError(err)
	fmt.Println(string(bodyBytes))

	resp.Body.Close()
	//c (api *Twitter) GetUserFollowers(id string, v url.Values, options ...QueueOption) (chan *Data, chan error)
	//resChan, errChan := s.tw.V2Client.GetUserFollowers("OpenAI", nil, twitter2.WithAuto(false))

	// Handle channels and timeout
	//var response *twitter2.Data
	//timeout := 10 * time.Second
	//var data []*twitter2.Tweet
	//
	//select {
	//case res := <-resChan:
	//	response = res
	//	b, err := json.Marshal(response.Data)
	//	s.Require().NoError(err)
	//	err = json.Unmarshal(b, &data)
	//	s.Require().NoError(err)
	//case err := <-errChan:
	//	fmt.Println("Error:", err)
	//	s.Require().NoError(err)
	//case <-time.After(timeout):
	//	fmt.Println("Error: Request timed out.")
	//	return
	//}
	//s.Assert().NotEmpty(data)
	//for _, tweet := range data {
	//	fmt.Printf("TweetID %s: AuthorID %s, Text: %s \n", tweet.ID, tweet.AuthorID, tweet.Text)
	//}
}

// 1590980065185058816

func (s *TwitterTestSuite) TestReadUserMe() {

	//userID := "1520852589067194373"

	r := resty_base.GetBaseRestyClient("https://api.twitter.com", addToken(consumerKey, consumerSecret))

	resp, err := r.R().Get("/2/users/me")
	s.Assert().NoError(err)
	s.Assert().NotNil(resp)
	bodyBytes, err := io.ReadAll(resp.RawBody())
	s.Assert().NoError(err)
	fmt.Println(string(bodyBytes))

	//// Close the response body
	//resp.Body.Close()
	//// 1520852589067194373

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
