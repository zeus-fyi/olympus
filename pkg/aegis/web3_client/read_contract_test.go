package web3_client

import (
	"github.com/zeus-fyi/gochain/web3"
)

func (s *Web3ClientTestSuite) TestReadContract() {
	//ctx := context.Background()
	//b, err := s.GoerliWeb3.GetCurrentBalance(ctx)
	//
	//s.Require().Nil(err)
	//s.Assert().NotNil(b)

	abis, err := web3.ABIOpenFile("/Users/alex/go/Olympus/olympus/pkg/aegis/web3_client/contract_abis/erc20Abi.json")
	s.Require().Nil(err)
	s.Assert().NotEmpty(abis)

}
