package artemis_uniswap_pricing

import (
	"fmt"

	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/constants"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/utils"
)

func (s *UniswapPricingTestSuite) TestPoolAddressesFromPricingStruct() {
	pair := []accounts.Address{artemis_trading_constants.WETH9ContractAddressAccount, artemis_trading_constants.PepeContractAddrAccount}
	pd, err := NewUniswapPools(pair)
	s.Require().Nil(err)

	fmt.Println(pd.V2Pair.PairContractAddr)
	fmt.Println(pd.V3Pairs.LowFee.PoolAddress)
	fmt.Println(pd.V3Pairs.MediumFee.PoolAddress)
	fmt.Println(pd.V3Pairs.HighFee.PoolAddress)
}

func (s *UniswapPricingTestSuite) TestPoolAddresses() {
	factoryAddress := artemis_trading_constants.UniswapV3FactoryAddressAccount
	tokenA := uniswap_core_entities.NewToken(1, artemis_trading_constants.WETH9ContractAddressAccount, 0, "", "")
	tokenB := uniswap_core_entities.NewToken(1, artemis_trading_constants.PepeContractAddrAccount, 0, "", "")

	pa, err := utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeMedium, "")
	s.Require().Nil(err)
	fmt.Println(pa.Hex())
	pa, err = utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeHigh, "")
	s.Require().Nil(err)
	fmt.Println(pa.Hex())

	v2Pair := UniswapV2Pair{}
	err = v2Pair.PairForV2(tokenA.Address.Hex(), tokenB.Address.Hex())
	s.Require().Nil(err)
	fmt.Println(v2Pair.PairContractAddr)
}
