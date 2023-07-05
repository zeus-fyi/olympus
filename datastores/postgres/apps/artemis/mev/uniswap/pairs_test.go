package artemis_models_uniswap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/constants"
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
	pair := []accounts.Address{artemis_trading_constants.WETH9ContractAddressAccount, artemis_trading_constants.PepeContractAddrAccount}
	err := InsertStandardUniswapPairInfoFromPair(ctx, pair)
	s.Require().Nil(err)
}

func TestUniswapModelsTestSuite(t *testing.T) {
	suite.Run(t, new(UniswapModelsTestSuite))
}
