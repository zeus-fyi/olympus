package web3_client

import (
	"fmt"
	"math/big"

	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
)

func (s *Web3ClientTestSuite) TestPricingImpact() {
	reserve0, _ := new(big.Int).SetString("400000", 10)  // TokenB
	reserve1, _ := new(big.Int).SetString("1200000", 10) // TokenA
	token0Addr, token1Addr := StringsToAddresses(PepeContractAddr, WETH9ContractAddress)
	mockPairResp := uniswap_pricing.UniswapV2Pair{
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
	mockPairResp = uniswap_pricing.UniswapV2Pair{
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
