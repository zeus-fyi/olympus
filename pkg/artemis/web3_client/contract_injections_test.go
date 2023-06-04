package web3_client

import "github.com/ethereum/go-ethereum/common/hexutil"

func (s *Web3ClientTestSuite) TestRawDawgInjection() {
	s.LocalHardhatMainnetUser.MustInjectRawDawg()
	bal := hexutil.Big{}
	bigInt := bal.ToInt()
	bigInt.Set(Ether)
	bal = hexutil.Big(*bigInt)
	err := s.LocalHardhatMainnetUser.SetBalance(ctx, RawDawgAddr, bal)
	s.Require().Nil(err)

	rawDawgBal, err := s.LocalHardhatMainnetUser.GetBalance(ctx, RawDawgAddr, nil)
	s.Require().Nil(err)
	s.Require().Equal(Ether, rawDawgBal)
}
