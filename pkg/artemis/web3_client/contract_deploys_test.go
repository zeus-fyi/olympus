package web3_client

import (
	"fmt"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

func (s *Web3ClientTestSuite) TestDeployRawdawgContract() {
	rawDawgPayload, bc := artemis_oly_contract_abis.MustLoadRawdawgContractDeployPayload()

	newAccount, err := accounts.ParsePrivateKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	s.Assert().Nil(err)
	s.LocalHardhatMainnetUser.Account = newAccount
	rawDawgPayload.GasLimit = 2000000
	rawDawgPayload.Params = []interface{}{}

	tx, err := s.LocalHardhatMainnetUser.DeploySmartContract(ctx, bc, rawDawgPayload)
	s.Require().Nil(err)
	s.Assert().NotNil(tx)

	rx, err := s.LocalHardhatMainnetUser.WaitForReceipt(ctx, tx.Hash())
	s.Assert().Nil(err)
	s.Assert().NotNil(rx)
}

func (s *Web3ClientTestSuite) TestDeployContract() {
	forceDirToLocation()
	tokenPayload, bc, err := artemis_oly_contract_abis.LoadERC20AbiPayload()
	s.Require().Nil(err)
	newAccount, err := accounts.ParsePrivateKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	s.Assert().Nil(err)
	s.LocalHardhatMainnetUser.Account = newAccount
	tokenPayload.GasLimit = 2000000

	mintAmount := new(big.Int).Mul(big.NewInt(10000000), Ether)
	tokenPayload.Params = []interface{}{mintAmount}

	tx, err := s.LocalHardhatMainnetUser.DeploySmartContract(ctx, bc, tokenPayload)
	s.Require().Nil(err)
	s.Assert().NotNil(tx)

	rx, err := s.LocalHardhatMainnetUser.WaitForReceipt(ctx, tx.Hash())
	s.Assert().Nil(err)
	s.Assert().NotNil(rx)

	b, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, rx.ContractAddress.String(), s.LocalHardhatMainnetUser.PublicKey())
	s.Require().Nil(err)
	s.Assert().NotZero(b)
	s.Assert().Equal(mintAmount.String(), b.String())

	b, err = s.LocalHardhatMainnetUser.GetBalance(ctx, s.LocalHardhatMainnetUser.PublicKey(), nil)
	s.Require().Nil(err)
	s.Assert().NotZero(b)

	s.Require().Nil(err)
	s.Assert().NotNil(tx)
}

func (s *Web3ClientTestSuite) TestDeployUniswapFactoryContract() {
	forceDirToLocation()
	factoryPayload, bc, err := artemis_oly_contract_abis.LoadUniswapFactoryAbiPayload()
	s.Require().Nil(err)
	newAccount, err := accounts.ParsePrivateKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	s.Assert().Nil(err)
	s.LocalHardhatMainnetUser.Account = newAccount
	factoryPayload.GasLimit = 2400000
	factoryPayload.Params = []interface{}{newAccount.PublicKey()}

	tx, err := s.LocalHardhatMainnetUser.DeployContract(ctx, bc, factoryPayload)
	s.Require().Nil(err)
	s.Assert().NotNil(tx)

	rx, err := s.LocalHardhatMainnetUser.WaitForReceipt(ctx, tx.Hash())
	s.Assert().Nil(err)
	s.Assert().NotNil(rx)
	fmt.Println(rx.ContractAddress.String())
}
