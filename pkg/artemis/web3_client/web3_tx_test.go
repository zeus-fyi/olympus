package web3_client

import (
	"fmt"

	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

func (s *Web3ClientTestSuite) TestWeb3SendEther() {
	b, err := s.GoerliWeb3User.GetCurrentBalance(ctx)
	s.Require().Nil(err)
	s.Assert().NotNil(b)

	params := web3_actions.SendEtherPayload{
		TransferArgs: web3_actions.TransferArgs{
			Amount:    Finney,
			ToAddress: s.GoerliWeb3User2.Address(),
		},
		GasPriceLimits: web3_actions.GasPriceLimits{},
	}
	send, err := s.GoerliWeb3User.Send(ctx, params)
	s.Require().Nil(err)
	s.Assert().NotNil(send)
}

func (s *Web3ClientTestSuite) TestWeb3TransferTokenToUser() {
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: LinkGoerliContractAddr,
		ContractFile:      web3_actions.ERC20,
		MethodName:        web3_actions.Transfer,
		SendEtherPayload: web3_actions.SendEtherPayload{
			GasPriceLimits: web3_actions.GasPriceLimits{},
		},
		Params: []interface{}{s.GoerliWeb3User2.Address(), Finney},
	}
	err := s.GoerliWeb3User.TransferERC20Token(ctx, params, false, 60)
	s.Require().Nil(err)
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
	fmt.Println(tx.Hash)

}
