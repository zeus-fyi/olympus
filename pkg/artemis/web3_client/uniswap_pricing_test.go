package web3_client

import (
	"fmt"
	"math/big"
)

func (s *Web3ClientTestSuite) TestPricingImpact() {
	reserve0, _ := new(big.Int).SetString("400000", 10)
	reserve1, _ := new(big.Int).SetString("1200000", 10)
	token0Addr, token1Addr := StringsToAddresses(PepeContractAddr, WETH9ContractAddress)
	mockPairResp := UniswapV2Pair{
		KLast:    big.NewInt(0),
		Token0:   token0Addr,
		Token1:   token1Addr,
		Reserve0: reserve0,
		Reserve1: reserve1,
	}
	originalRate, _ := mockPairResp.GetToken1Price()
	newRate := mockPairResp.PriceImpact(big.NewInt(3000))
	fmt.Println("originalRate", originalRate)
	fmt.Println("newRate", newRate)
}
