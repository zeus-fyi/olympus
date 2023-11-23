package hera_discord

import (
	"context"
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
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	token, err := read_keys.GetDiscordKey(ctx, s.Tc.ProductionLocalTemporalUserID)
	s.Require().Nil(err)
	s.Require().NotEmpty(token)
	InitDiscordClient(ctx, token)
	s.Require().NotNil(DiscordClient.DC)

}

// https://discord.com/api/oauth2/authorize?client_id=1177043826133717112&redirect_uri=http%3A%2F%2Flocalhost%3A9002%2Fdiscord%2Fcallback&response_type=code&scope=guilds%20messages.read%20guilds.join%20guilds.members.read
func (s *DiscordTestSuite) TestReadPosts() {

	d := DiscordClient
	r, err := d.ListAllChannels(ctx)
	s.Require().Nil(err)
	s.Require().NotEmpty(r)
}

func TestDiscordTestSuite(t *testing.T) {
	suite.Run(t, new(DiscordTestSuite))
}
