package web3_client

import (
	"fmt"
	"math/big"
)

func (s *Web3ClientTestSuite) TestPricingImpact() {
	reserve0, _ := new(big.Int).SetString("400000", 10)  // TokenB
	reserve1, _ := new(big.Int).SetString("1200000", 10) // TokenA
	token0Addr, token1Addr := StringsToAddresses(PepeContractAddr, WETH9ContractAddress)
	mockPairResp := UniswapV2Pair{
		KLast:    big.NewInt(0),
		Token0:   token0Addr,
		Token1:   token1Addr,
		Reserve0: reserve0,
		Reserve1: reserve1,
	}
	originalRate, _ := mockPairResp.GetToken1Price()
	newRateToken1, newRateToken0 := mockPairResp.PriceImpactToken1BuyToken0(big.NewInt(3000))
	fmt.Println("originalRate", originalRate)
	fmt.Println("newRateToken0", newRateToken0)
	fmt.Println("newRateToken1", newRateToken1)
	s.Assert().Equal("3.015037481", newRateToken0.String())

	reserve0, _ = new(big.Int).SetString("400000", 10)  // TokenB
	reserve1, _ = new(big.Int).SetString("1200000", 10) // TokenA
	mockPairResp = UniswapV2Pair{
		KLast:    big.NewInt(0),
		Token0:   token0Addr,
		Token1:   token1Addr,
		Reserve0: reserve0,
		Reserve1: reserve1,
	}
	originalRate, _ = mockPairResp.GetToken1Price()
	newRateToken1, newRateToken0 = mockPairResp.PriceImpactToken0BuyToken1(big.NewInt(1000))
	fmt.Println("originalRate", originalRate)
	fmt.Println("newRateToken0", newRateToken0)
	fmt.Println("newRateToken1", newRateToken1)
	s.Assert().Equal("2.985037518", newRateToken0.String())
}
