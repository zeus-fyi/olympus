package web3_client

import (
	"fmt"

	"github.com/gochain/gochain/v4/common/hexutil"
)

func (s *Web3ClientTestSuite) TestSetBalance() {
	cb, err := s.LocalMainnetWeb3User.GetCurrentBalance(ctx)
	fmt.Println("startingBalance", cb.String())
	bal := hexutil.Big{}
	bigInt := bal.ToInt()
	bigInt.Set(Ether)
	bal = hexutil.Big(*bigInt)
	err = s.LocalHardhatMainnetUser.SetBalance(ctx, s.LocalMainnetWeb3User.PublicKey(), bal)
	s.Require().Nil(err)
	cb, err = s.LocalMainnetWeb3User.GetCurrentBalance(ctx)
	s.Require().Nil(err)
	s.Assert().Equal(Ether.String(), cb.String())
}

func (s *Web3ClientTestSuite) TestResetNetwork() {
	err := s.LocalMainnetWeb3User.ResetNetwork(ctx, s.Tc.HardhatNode, 17317757)
	s.Require().Nil(err)
}
