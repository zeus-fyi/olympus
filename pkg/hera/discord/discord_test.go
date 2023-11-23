package hera_discord

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type DiscordTestSuite struct {
	test_suites_base.TestSuite
}

func (s *DiscordTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	token, err := read_keys.GetDiscordKey(ctx, s.Tc.ProductionLocalTemporalUserID)
	s.Require().Nil(err)
	s.Require().NotEmpty(token)
	InitDiscordClient(ctx, token)
	s.Require().NotNil(DiscordClient.Client)
}

func (s *DiscordTestSuite) TestReadPosts() {
	channels, err := DiscordClient.ListAllChannels(ctx)
	s.Require().Nil(err)
	s.Require().NotNil(channels)
	for _, chn := range channels {
		fmt.Println(chn.Name, chn.ID)
	}
}

func TestDiscordTestSuite(t *testing.T) {
	suite.Run(t, new(DiscordTestSuite))
}
