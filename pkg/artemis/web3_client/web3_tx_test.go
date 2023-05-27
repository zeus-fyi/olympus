package web3_client

import (
	"fmt"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (s *Web3ClientTestSuite) TestWeb3SendEther() {
	b, err := s.GoerliWeb3User.GetCurrentBalance(ctx)
	s.Require().Nil(err)
	s.Assert().NotNil(b)

	params := web3_actions.SendEtherPayload{
		TransferArgs: web3_actions.TransferArgs{
			Amount:    Finney,
			ToAddress: accounts.Address(s.GoerliWeb3User2.Address()),
		},
		GasPriceLimits: web3_actions.GasPriceLimits{},
	}
	send, err := s.GoerliWeb3User.Send(ctx, params)
	s.Require().Nil(err)
	s.Assert().NotNil(send)
}

func (s *Web3ClientTestSuite) TestWeb3TransferTokenToUser() {
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: LinkTokenAddr,
		ContractFile:      web3_actions.ERC20,
		MethodName:        web3_actions.Transfer,
		SendEtherPayload: web3_actions.SendEtherPayload{
			GasPriceLimits: web3_actions.GasPriceLimits{},
		},
		Params: []interface{}{s.LocalHardhatMainnetUser.Address(), Finney},
	}
	tx, err := s.LocalHardhatMainnetUser.TransferERC20Token(ctx, params)
	s.Require().Nil(err)
	s.Require().NotNil(tx)
}

func (s *Web3ClientTestSuite) TestWeb3TransferTokenToUserFromPresignedTx() {
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: LinkGoerliContractAddr,
		ContractFile:      web3_actions.ERC20,
		MethodName:        web3_actions.Transfer,
		SendEtherPayload: web3_actions.SendEtherPayload{
			GasPriceLimits: web3_actions.GasPriceLimits{},
		},
		Params: []interface{}{s.GoerliWeb3User2.Address(), Finney},
	}
	signedTx, err := s.GoerliWeb3User.GetSignedTxToCallFunctionWithArgs(ctx, &params)
	s.Require().Nil(err)
	s.Require().NotNil(signedTx)

	tx, err := s.GoerliWeb3User.SubmitSignedTxAndReturnTxData(ctx, signedTx)
	s.Require().Nil(err)
	s.Require().NotEmpty(tx)
	fmt.Println(tx.Hash())
}
