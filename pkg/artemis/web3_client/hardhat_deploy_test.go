package web3_client

import (
	"fmt"

	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

// this test shows you can replace any address with any bytecode, this was an erc20 override and I forced the balance to be 1 ether
func (s *Web3ClientTestSuite) TestSetCodeOverride() {
	forceDirToLocation()

	randomAddr := "0x7623e9dc0da6ff821ddb9ebaba794054e078f8c4"

	bc, err := artemis_oly_contract_abis.LoadERC20DeployedByteCode()
	s.Require().NoError(err)
	s.Require().NotNil(bc)
	err = s.LocalHardhatMainnetUser.SetCodeOverride(ctx, randomAddr, bc)
	s.Require().Nil(err)
	newAccount, err := accounts.ParsePrivateKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	s.Assert().Nil(err)
	s.LocalHardhatMainnetUser.Account = newAccount
	s.Require().Nil(err)

	b, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, randomAddr, s.LocalHardhatMainnetUser.PublicKey())
	s.Require().Nil(err)
	s.Assert().NotZero(b)
	fmt.Println(b.String())

	err = s.LocalHardhatMainnetUser.SetERC20BalanceBruteForce(ctx, randomAddr, s.LocalHardhatMainnetUser.PublicKey(), Ether)
	s.Require().Nil(err)
	b, err = s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, randomAddr, s.LocalHardhatMainnetUser.PublicKey())
	s.Require().Nil(err)
	s.Assert().NotZero(b)
	fmt.Println(b.String())
}
