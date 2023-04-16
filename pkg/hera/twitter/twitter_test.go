package hera_twitter

import (
	"context"
	"testing"

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
	tw := InitTwitterClient(ctx, s.Tc.TwitterBearerToken)
	s.Assert().NotNil(tw)
	s.tw = tw
}

func (s *TwitterTestSuite) TestReadPosts() {

}

func TestTwitterTestSuite(t *testing.T) {
	suite.Run(t, new(TwitterTestSuite))
}
