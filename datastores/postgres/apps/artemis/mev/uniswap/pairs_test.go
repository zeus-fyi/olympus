package artemis_models_uniswap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

var ctx = context.Background()

type UniswapModelsTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *UniswapModelsTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
}

func (s *UniswapModelsTestSuite) TestInsertPair() {
	p := artemis_autogen_bases.UniswapPairInfo{
		TradingEnabled:    false,
		Address:           "",
		FactoryAddress:    "",
		Fee:               0,
		Version:           "",
		Token0:            "",
		Token1:            "",
		ProtocolNetworkID: 0,
	}
	err := InsertUniswapPairInfo(ctx, p)
	s.Require().Nil(err)
}

func TestUniswapModelsTestSuite(t *testing.T) {
	suite.Run(t, new(UniswapModelsTestSuite))
}
