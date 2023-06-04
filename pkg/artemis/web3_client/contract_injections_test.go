package web3_client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

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

	abiInfo := MustLoadRawdawgAbi()
	owner, err := s.LocalHardhatMainnetUser.GetOwner(ctx, abiInfo, RawDawgAddr)
	s.Require().Nil(err)
	fmt.Println(owner.String())
	// now try doing a swap
}
