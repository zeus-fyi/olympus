package web3_client

import (
	"fmt"
	"math/big"

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

func (s *Web3ClientTestSuite) TestAccountImpersonation() {
	userToImpersonate := "0x5c3fd6932ce20b60af632d8983c0121db7beef46"
	err := s.LocalMainnetWeb3User.ImpersonateAccount(ctx, userToImpersonate)
	s.Require().Nil(err)
	err = s.LocalMainnetWeb3User.StopImpersonatingAccount(ctx, userToImpersonate)
	s.Require().Nil(err)
	cb, err := s.LocalMainnetWeb3User.GetUserCurrentBalance(ctx, userToImpersonate)
	s.Require().Nil(err)
	s.Assert().NotZero(cb)
}

func (s *Web3ClientTestSuite) TestGetEvmSnapshot() {
	ss, err := s.LocalMainnetWeb3User.GetEvmSnapshot(ctx)
	s.Require().Nil(err)
	s.Assert().NotZero(ss)
}

func (s *Web3ClientTestSuite) TestGetSlot() {
	zeroAddr := "0x0000000000000000000000000000000000000000"
	slotNum := new(big.Int).SetUint64(0)
	hexStr, _ := getSlot(zeroAddr, slotNum)
	fmt.Println("hexStr", hexStr)

	s.Assert().Equal("0xad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5", hexStr)
}

func (s *Web3ClientTestSuite) TestGetSlotFromKnownNonZeroERC20Balance() {
	// block set to 17317757
	usdtAddr := "0xdac17f958d2ee523a2206206994597c13d831ec7"
	userAddr := "0xC6CDE7C39eB2f0F0095F41570af89eFC2C1Ea828"
	b, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, usdtAddr, userAddr)
	s.Require().Nil(err)
	_, slotHex, err := s.LocalHardhatMainnetUser.FindSlotFromUserWithBalance(ctx, usdtAddr, userAddr)
	s.Require().Nil(err)
	fmt.Println("slotHex", slotHex)

	resp, err := s.LocalHardhatMainnetUser.GetStorageAt(ctx, usdtAddr, slotHex)
	s.Require().Nil(err)
	foundBal := new(big.Int).SetBytes(resp)
	s.Assert().Equal(b.String(), foundBal.String())
}

func (s *Web3ClientTestSuite) TestSetERC20BalanceAtSlotNumber() {
	// block set to 17317757
	usdtAddr := "0xdac17f958d2ee523a2206206994597c13d831ec7"
	b, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, usdtAddr, s.LocalHardhatMainnetUser.PublicKey())
	s.Require().Nil(err)
	fmt.Println("startingBalance", b.String())
	err = s.LocalHardhatMainnetUser.SetERC20BalanceAtSlotNumber(ctx, usdtAddr, s.LocalMainnetWeb3User.PublicKey(), 2, Ether)
	s.Require().Nil(err)
	b, err = s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, usdtAddr, s.LocalHardhatMainnetUser.PublicKey())
	s.Require().Nil(err)
	s.Assert().Equal(Ether.String(), b.String())
}

func (s *Web3ClientTestSuite) TestSetERC20BalanceBruteForce() {
	// block set to 17317757
	err := s.LocalMainnetWeb3User.ResetNetwork(ctx, s.Tc.HardhatNode, 17317757)
	s.Require().Nil(err)
	usdtAddr := "0xdac17f958d2ee523a2206206994597c13d831ec7"
	b, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, usdtAddr, s.LocalHardhatMainnetUser.PublicKey())
	s.Require().Nil(err)
	fmt.Println("startingBalance", b.String())
	err = s.LocalHardhatMainnetUser.SetERC20BalanceBruteForce(ctx, usdtAddr, s.LocalHardhatMainnetUser.PublicKey(), Ether)
	s.Require().Nil(err)
	b, err = s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, usdtAddr, s.LocalHardhatMainnetUser.PublicKey())
	s.Require().Nil(err)
	s.Assert().Equal(Ether.String(), b.String())
	fmt.Println("endingBalance", b.String())
}
