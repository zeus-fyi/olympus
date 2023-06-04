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

	// DAI-USDC pair contract
	daiUsdcPairContractAddr := "0xAE461cA67B15dc8dc81CE7615e0320dA1A9aB8D5"
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
	pair, err := uni.GetPairContractPrices(ctx, daiUsdcPairContractAddr)
	s.Require().Nil(err)
	s.Require().NotEmpty(pair)

	// needs to give contract tokens

	// now try doing a swap
	// just fork mainnet and try to swap a common token pair

	/*
	   address _pair,
	   address _token_in,
	   uint256 _amountIn,
	   uint256 _amountOut,
	   bool _isToken0
	*/
}
