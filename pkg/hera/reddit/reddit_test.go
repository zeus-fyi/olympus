package hera_reddit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type RedditTestSuite struct {
	test_suites_base.TestSuite
	rc Reddit
}

func (s *RedditTestSuite) SetupTest() {
	s.InitLocalConfigs()
	rc, err := InitRedditClient(ctx, s.Tc.RedditPublicOAuth2, s.Tc.RedditSecretOAuth2, s.Tc.RedditUsername, s.Tc.RedditPassword)
	s.Require().Nil(err)
	s.Assert().NotNil(rc)
	s.rc = rc
}

func (s *RedditTestSuite) TestReadPosts() {
	posts, _, err := s.rc.ReadOnly.Subreddit.TopPosts(ctx, "kubernetes", &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 1,
		},
		Time: "day",
	})
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	searchID := 1700685362343578112

	resp, err := hera_search.InsertIncomingRedditPosts(ctx, searchID, posts)
	s.Require().Nil(err)
	s.Assert().NotZero(resp)
}

func TestRedditTestSuite(t *testing.T) {
	suite.Run(t, new(RedditTestSuite))
}
