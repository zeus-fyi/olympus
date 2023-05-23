package web3_client

import (
	"github.com/gochain/gochain/v4/common/hexutil"
)

func (s *Web3ClientTestSuite) TestSetBalance() {
	bal := hexutil.Big{}
	bigInt := bal.ToInt()
	bigInt.Set(Ether)
	bal = hexutil.Big(*bigInt)
	err := s.LocalHardhatMainnetUser.SetBalance(ctx, s.LocalMainnetWeb3User.PublicKey(), bal)
	s.Require().Nil(err)
	cb, err := s.LocalMainnetWeb3User.GetCurrentBalance(ctx)
	s.Require().Nil(err)
	s.Assert().Equal(Ether.String(), cb.String())
}
