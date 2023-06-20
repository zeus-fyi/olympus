package web3_client

import (
	"github.com/zeus-fyi/gochain/web3/accounts"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/utils"
)

func (s *Web3ClientTestSuite) TestUniswapQuoterV2() {
	factoryAddress := accounts.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")
	tokenA := core_entities.NewToken(1, accounts.HexToAddress(WETH9ContractAddress), 18, "WETH", "Wrapped Ether")
	tokenB := core_entities.NewToken(1, accounts.HexToAddress(LooksTokenAddr), 18, "LOOKS", "LooksRare Token")

	result, err := utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeMedium, "")
	s.Require().Nil(err)
	s.Equal(result.String(), accounts.HexToAddress("0x4b5Ab61593A2401B1075b90c04cBCDD3F87CE011").String())

	v3Pool, err := entities.NewPool(tokenA, tokenB, constants.FeeLow, nil, nil, 0, nil)
	s.Require().NoError(err)
	s.Require().NotNil(v3Pool)
	//
	//uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
	//resp, err := uni.GetPoolV3QuoteFromQuoterV2(ctx, entities.Pool{})
	//s.Require().NoError(err)
	//s.Require().NotNil(resp)
}
