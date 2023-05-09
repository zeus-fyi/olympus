package web3_client

import "github.com/gochain/gochain/v4/common"

func (s *Web3ClientTestSuite) TestGetPairContract() {
	uni := InitUniswapV2Client(ctx)

	addr := uni.GetPairContract(common.HexToAddress(""), common.HexToAddress(""))
	s.Assert().Equal(common.HexToAddress(""), addr)
}
