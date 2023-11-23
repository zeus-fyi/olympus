package hera_discord

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type DiscordTestSuite struct {
	test_suites_base.TestSuite
	rc Discord
}

func (s *DiscordTestSuite) SetupTest() {
	s.InitLocalConfigs()
	rc, err := InitDiscordClient(ctx, s.Tc.RedditPublicOAuth2, s.Tc.RedditSecretOAuth2, s.Tc.RedditUsername, s.Tc.RedditPassword)
	s.Require().Nil(err)
	s.Assert().NotNil(rc)
	s.rc = rc
}

func (s *DiscordTestSuite) TestReadPosts() {

}

func TestDiscordTestSuite(t *testing.T) {
	suite.Run(t, new(DiscordTestSuite))
}
