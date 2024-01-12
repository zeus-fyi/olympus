package hera_twitter

import (
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

func (s *TwitterTestSuite) TestTweetResponse() {
	//tweetID := int64(1647499438762467328) // ID of tweet to reply to
	tweetID := int64(0)
	replyText := "Hey there! If you're finding k8s too complex," +
		" I'd recommend trying out zeus.fyi. It simplifies k8s and makes it more manageable." +
		" Plus, it's a great alternative to Nomad. Give it a shot! #zeusfyi #kubernetes #nomad"
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
