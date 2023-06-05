package web3_client

import (
	"fmt"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (s *Web3ClientTestSuite) TestMintTokens() {
	forceDirToLocation()
	tokenPayload, bc, err := LoadERC20AbiPayload()
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

	swapContractAddr := "0x5FbDB2315678afecb367f032d93F642f64180aa3"
	_, err = s.LocalHardhatMainnetUser.ERC20ApproveSpender(ctx, rx.ContractAddress.String(), swapContractAddr, Ether)
	s.Require().Nil(err)

	transferTxParams := web3_actions.SendContractTxPayload{
		SmartContractAddr: rx.ContractAddress.String(),
		SendEtherPayload: web3_actions.SendEtherPayload{
			TransferArgs: web3_actions.TransferArgs{
				ToAddress: accounts.HexToAddress(swapContractAddr),
			},
		},
		ContractABI: MustLoadERC20Abi(),
		Params:      []interface{}{accounts.HexToAddress(swapContractAddr), Ether},
	}
	_, err = s.LocalHardhatMainnetUser.TransferERC20Token(ctx, transferTxParams)
	s.Require().Nil(err)

	b, err = s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, rx.ContractAddress.String(), swapContractAddr)
	s.Require().Nil(err)
	s.Assert().NotZero(b)
	fmt.Println(b.String())
}
