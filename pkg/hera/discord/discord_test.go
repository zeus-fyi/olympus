package hera_discord

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

var ctx = context.Background()

type DiscordTestSuite struct {
	test_suites_base.TestSuite
}

func (s *DiscordTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	authToken, err := read_keys.GetDiscordKey(ctx, s.Tc.ProductionLocalTemporalUserID)
	s.Require().Nil(err)
	s.Require().NotEmpty(authToken)
}

func (s *DiscordTestSuite) TestFetchChatMessages() {
	f := filepaths.Path{
		DirIn: "/Users/alex/go/Olympus/olympus/pkg/hera/discord",
		FnIn:  "qn.json",
	}

	b := f.ReadFileInPath()
	s.Require().NotNil(b)

	zzz := ChannelMessages{}
	err := json.Unmarshal(b, &zzz)
	s.Require().Nil(err)
	s.Require().NotEmpty(zzz)

	//for _, cm := range zzz.Messages {
	//	fmt.Println(cm.TimestampEdited.String())
	//}
}

func TestDiscordTestSuite(t *testing.T) {
	suite.Run(t, new(DiscordTestSuite))
}
