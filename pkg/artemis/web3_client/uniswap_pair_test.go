package web3_client

import (
	"fmt"
)

func (s *Web3ClientTestSuite) TestGetPairContract() {
	uni := InitUniswapV2Client(ctx, s.MainnetWeb3User)
	pair := uni.GetPairContractFromFactory(ctx, WETH9ContractAddress, PepeContractAddr)
	s.Assert().NotEmpty(pair)
	fmt.Println(pair.String())
	s.Assert().Equal("0xA43fe16908251ee70EF74718545e4FE6C5cCEc9f", pair.String())

	pair = uni.GetPairContractFromFactory(ctx, PepeContractAddr, WETH9ContractAddress)
}

func (s *Web3ClientTestSuite) TestGetPairContractInfo() {
	uni := InitUniswapV2Client(ctx, s.MainnetWeb3User)
	pairContractAddr := "0xA43fe16908251ee70EF74718545e4FE6C5cCEc9f"
	pair, err := uni.GetPairContractPrices(ctx, pairContractAddr)
	s.Assert().Nil(err)
	s.Assert().NotEmpty(pair)
	fmt.Println(pair.Token0.String())
	fmt.Println(pair.Token1.String())
	s.Assert().Equal(PepeContractAddr, pair.Token0.String())
	s.Assert().Equal(WETH9ContractAddress, pair.Token1.String())
	fmt.Println(pair.KLast.String())
	fmt.Println(pair.Reserve0.Uint64())
	fmt.Println(pair.Reserve1.Uint64())
}
