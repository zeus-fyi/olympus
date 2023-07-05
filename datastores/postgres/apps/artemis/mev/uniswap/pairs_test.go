package artemis_models_uniswap

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/constants"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/utils"
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
	//p := artemis_autogen_bases.UniswapPairInfo{
	//	TradingEnabled:    false,
	//	Address:           "",
	//	FactoryAddress:    "",
	//	Fee:               0,
	//	Version:           "",
	//	Token0:            "",
	//	Token1:            "",
	//	ProtocolNetworkID: 0,
	//}
	//err := InsertUniswapPairInfo(ctx, p)
	//s.Require().Nil(err)

	factoryAddress := artemis_trading_constants.UniswapV3FactoryAddressAccount
	tokenA := uniswap_core_entities.NewToken(1, artemis_trading_constants.WETH9ContractAddressAccount, 0, "", "")
	tokenB := uniswap_core_entities.NewToken(1, artemis_trading_constants.PepeContractAddrAccount, 0, "", "")

	pa, err := utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeMedium, "")
	s.Require().Nil(err)
	fmt.Println(pa.Hex())
	pa, err = utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeHigh, "")
	s.Require().Nil(err)
	fmt.Println(pa.Hex())

	v2Pair := uniswap_pricing.UniswapV2Pair{}
	err = v2Pair.PairForV2(tokenA.Address.Hex(), tokenB.Address.Hex())
	s.Require().Nil(err)
	fmt.Println(v2Pair.PairContractAddr)
}

func TestUniswapModelsTestSuite(t *testing.T) {
	suite.Run(t, new(UniswapModelsTestSuite))
}
