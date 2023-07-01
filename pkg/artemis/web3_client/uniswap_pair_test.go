package web3_client

import (
	"fmt"
)

func (s *Web3ClientTestSuite) TestGetPairContract() {
	uni := InitUniswapClient(ctx, s.MainnetWeb3User)
	pair := uni.GetPairContractFromFactory(ctx, WETH9ContractAddress, PepeContractAddr)
	s.Assert().NotEmpty(pair)
	fmt.Println(pair.String())
	s.Assert().Equal("0xA43fe16908251ee70EF74718545e4FE6C5cCEc9f", pair.String())

	pair = uni.GetPairContractFromFactory(ctx, PepeContractAddr, WETH9ContractAddress)
}

func (s *Web3ClientTestSuite) TestGetPairContractInfo() {
	//uni := InitUniswapClient(ctx, s.MainnetWeb3User)
	//pairContractAddr := "0xA43fe16908251ee70EF74718545e4FE6C5cCEc9f"
	//pair, err := uni.GetPairContractPrices(ctx, pairContractAddr)
	//s.Assert().Nil(err)
	//s.Assert().NotEmpty(pair)
	//fmt.Println("token0", pair.Token0.String())
	//fmt.Println("token1", pair.Token1.String())
	//s.Assert().Equal(PepeContractAddr, pair.Token0.String())
	//s.Assert().Equal(WETH9ContractAddress, pair.Token1.String())
	//fmt.Println("kLast", pair.KLast.String())
	//fmt.Println("reserve0", pair.Reserve0.String())
	//fmt.Println("reserve1", pair.Reserve1.String())
	//fmt.Println("price0CumulativeLast", pair.Price0CumulativeLast.String())
	//fmt.Println("price1CumulativeLast", pair.Price1CumulativeLast.String())
}
