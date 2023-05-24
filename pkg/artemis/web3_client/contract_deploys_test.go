package web3_client

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

func (s *Web3ClientTestSuite) TestDeployContract() {
	forceDirToLocation()
	tokenPayload, bc, err := LoadERC20AbiPayload()
	s.Require().Nil(err)
	newAccount, err := accounts.ParsePrivateKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	s.Assert().Nil(err)
	s.LocalHardhatMainnetUser.Account = newAccount
	tokenPayload.GasLimit = 2000000

	mintAmount := new(big.Int).Mul(big.NewInt(10000000), Ether)
	tokenPayload.Params = []interface{}{mintAmount}

	tx, err := s.LocalHardhatMainnetUser.DeployERC20Token(ctx, bc, tokenPayload)
	s.Require().Nil(err)
	s.Assert().NotNil(tx)

	rx, err := s.LocalHardhatMainnetUser.WaitForReceipt(ctx, tx.Hash)
	s.Assert().Nil(err)
	s.Assert().NotNil(rx)

	b, err := s.LocalMainnetWeb3User.ReadERC20TokenBalance(ctx, rx.ContractAddress.String(), s.LocalHardhatMainnetUser.PublicKey())
	s.Require().Nil(err)
	s.Assert().NotZero(b)
	s.Assert().Equal(mintAmount.String(), b.String())
}
