package web3_client

import (
	"github.com/zeus-fyi/gochain/web3/types"
)

func (s *Web3ClientTestSuite) TestReadContract() {
	//ctx := context.Background()
	//b, err := s.GoerliWeb3User.GetCurrentBalance(ctx)
	//
	//s.Require().Nil(err)
	//s.Assert().NotNil(b)

	abis, err := web3_types.ABIOpenFile("/Users/alex/go/Olympus/olympus/pkg/aegis/web3_client/contract_abis/erc20Abi.json")
	s.Require().Nil(err)
	s.Assert().NotEmpty(abis)

}
