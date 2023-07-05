package artemis_models_uniswap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/constants"
)

var ctx = context.Background()

type UniswapModelsTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *UniswapModelsTestSuite) SetupTest() {
	s.InitLocalConfigs()
	//s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	s.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}

func (s *UniswapModelsTestSuite) TestInsertPair() {

	tokens, _, err := artemis_mev_models.SelectERC20Tokens(ctx)
	s.Require().Nil(err)
	var tokensToAdd []accounts.Address
	for _, to := range tokens {
		addressStr := to.Address
		tokenAddr := accounts.HexToAddress(addressStr)
		if tokenAddr != artemis_trading_constants.WETH9ContractAddressAccount {
			tokensToAdd = append(tokensToAdd, tokenAddr)
		}
	}

	for _, toAdd := range tokensToAdd {
		pair := []accounts.Address{artemis_trading_constants.WETH9ContractAddressAccount, toAdd}
		err = InsertStandardUniswapPairInfoFromPair(ctx, pair)
		s.Require().Nil(err)
	}
}

func TestUniswapModelsTestSuite(t *testing.T) {
	suite.Run(t, new(UniswapModelsTestSuite))
}
