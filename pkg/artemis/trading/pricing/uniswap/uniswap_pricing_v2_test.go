package artemis_uniswap_pricing

import (
	"fmt"
	"math/big"

	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_utils "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/utils"
)

func (s *UniswapPricingTestSuite) TestPricingImpact() {
	reserve0, _ := new(big.Int).SetString("400000", 10)  // TokenB
	reserve1, _ := new(big.Int).SetString("1200000", 10) // TokenA
	token0Addr, token1Addr := artemis_utils.StringsToAddresses(artemis_trading_constants.PepeContractAddr, WETH9ContractAddress)
	mockPairResp := UniswapV2Pair{
		KLast:    big.NewInt(0),
		Token0:   token0Addr,
		Token1:   token1Addr,
		Reserve0: reserve0,
		Reserve1: reserve1,
	}
	to, reservesToken0, reservesToken1 := mockPairResp.PriceImpactToken1BuyToken0(big.NewInt(3000))
	fmt.Println("to.AmountOut", to.AmountOut.String())
	fmt.Println("reservesToken0", reservesToken0.String())
	fmt.Println("reservesToken1", reservesToken1.String())
	s.Assert().Equal(big.NewInt(399006), reservesToken0)
	s.Assert().Equal(big.NewInt(1203000), reservesToken1)

	reserve0, _ = new(big.Int).SetString("400000", 10)  // TokenB
	reserve1, _ = new(big.Int).SetString("1200000", 10) // TokenA
	mockPairResp = UniswapV2Pair{
		KLast:    big.NewInt(0),
		Token0:   token0Addr,
		Token1:   token1Addr,
		Reserve0: reserve0,
		Reserve1: reserve1,
	}
	to, reservesToken0, reservesToken1 = mockPairResp.PriceImpactToken0BuyToken1(big.NewInt(1000))
	fmt.Println("to.AmountOut", to.AmountOut.String())
	fmt.Println("reservesToken0", reservesToken0.String())
	fmt.Println("reservesToken1", reservesToken1.String())
}

func (s *UniswapPricingTestSuite) TestUniswapSortTokens() {
	p := UniswapV2Pair{}
	err := p.PairForV2(artemis_trading_constants.PepeContractAddr, WETH9ContractAddress)
	s.Require().Nil(err)
	s.Require().Equal("0xA43fe16908251ee70EF74718545e4FE6C5cCEc9f", p.PairContractAddr)
	s.Require().Equal(p.Token0.String(), artemis_trading_constants.PepeContractAddr)
	s.Require().Equal(p.Token1.String(), WETH9ContractAddress)
}
