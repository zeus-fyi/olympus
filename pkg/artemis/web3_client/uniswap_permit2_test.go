package web3_client

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

func (s *Web3ClientTestSuite) TestPermit2PermitBatchEncode() {
	addr1 := accounts.HexToAddress(LidoSEthAddr)
	addr2 := accounts.HexToAddress(WETH9ContractAddress)
	permit2Batch := Permit2PermitBatchParams{
		PermitBatch: PermitBatch{
			Details: []PermitDetails{{
				Token:      addr1,
				Amount:     new(big.Int).SetUint64(1000000000000000000),
				Expiration: new(big.Int).SetUint64(1000000000000000000),
				Nonce:      new(big.Int).SetUint64(1),
			}},
			Spender:     addr2,
			SigDeadline: new(big.Int).SetUint64(1000000000000000000),
		},
		Signature: []byte("test"),
	}
	// convert to command
	ur := UniversalRouterExecCmd{
		Commands: []UniversalRouterExecSubCmd{
			{
				Command:       Permit2PermitBatch,
				CanRevert:     true,
				Inputs:        nil,
				DecodedInputs: permit2Batch,
			},
		},
	}
	encCmd, err := ur.EncodeCommands(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(encCmd)
	s.Require().NotNil(encCmd.Commands)
	subCmd := UniversalRouterExecSubCmd{}
	for i, byteVal := range encCmd.Commands {
		err = subCmd.DecodeCommand(byteVal, encCmd.Inputs[i])
		s.Require().NoError(err)
		s.Assert().Equal(true, subCmd.CanRevert)
		s.Assert().Equal(Permit2PermitBatch, subCmd.Command)
		decodedInputs := subCmd.DecodedInputs.(Permit2PermitBatchParams)
		s.Assert().Equal(permit2Batch.Signature, decodedInputs.Signature)
	}
}

func (s *Web3ClientTestSuite) TestPermit2PermitBatchEncode2() {
	addr1 := accounts.HexToAddress(LidoSEthAddr)
	//addr2 := accounts.HexToAddress(WETH9ContractAddress)
	permit2TransferFromBatch := Permit2PermitTransferFromBatchParams{
		Details: []AllowanceTransferDetails{
			{
				From:   s.LocalMainnetWeb3User.Address(),
				To:     accounts.HexToAddress(UniswapUniversalRouterAddress),
				Amount: new(big.Int).SetUint64(1000000000000000000),
				Token:  addr1,
			},
		},
	}
	// convert to command
	ur := UniversalRouterExecCmd{
		Commands: []UniversalRouterExecSubCmd{
			{
				Command:       Permit2TransferFromBatch,
				CanRevert:     true,
				Inputs:        nil,
				DecodedInputs: permit2TransferFromBatch,
			},
		},
	}
	encCmd, err := ur.EncodeCommands(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(encCmd)
	s.Require().NotNil(encCmd.Commands)
	subCmd := UniversalRouterExecSubCmd{}
	for i, byteVal := range encCmd.Commands {
		err = subCmd.DecodeCommand(byteVal, encCmd.Inputs[i])
		s.Require().NoError(err)
		s.Assert().Equal(true, subCmd.CanRevert)
		s.Assert().Equal(Permit2TransferFromBatch, subCmd.Command)
		decodedInputs := subCmd.DecodedInputs.(Permit2PermitTransferFromBatchParams)
		s.Assert().Equal(permit2TransferFromBatch.Details[0].From, decodedInputs.Details[0].From)
	}
}
