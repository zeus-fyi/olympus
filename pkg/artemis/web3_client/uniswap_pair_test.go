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
}
