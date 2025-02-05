package web3_client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

var (
	chainID = big.NewInt(1)
	urAddr  = accounts.HexToAddress(UniswapUniversalRouterAddressNew)
)

func (s *Web3ClientTestSuite) TestCopyPermitTest1() {
	expiration, _ := new(big.Int).SetString("3000000000000", 10)
	sigDeadline, _ := new(big.Int).SetString("3000000000000", 10)
	amount := new(big.Int).SetUint64(1000000000)
	usdcAddr := accounts.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	pp := Permit2PermitParams{
		PermitSingle: PermitSingle{
			PermitDetails: PermitDetails{
				Token:      usdcAddr,
				Amount:     amount,
				Expiration: expiration,
				Nonce:      new(big.Int).SetUint64(0),
			},
			Spender:     accounts.HexToAddress("0xe808c1cfeebb6cb36b537b82fa7c9eef31415a05"),
			SigDeadline: sigDeadline,
		},
		Signature: nil,
	}

	permitAddress := "0x4a873bdd49f7f9cc0a5458416a12973fab208f8d"
	err := pp.Sign(s.LocalHardhatMainnetUser.Account, chainID, accounts.HexToAddress(permitAddress), "Permit2")
	s.Require().Nil(err)
	s.Require().NotNil(pp.Signature)

	hashed := hashPermitSingle(pp.PermitSingle)
	eip := NewEIP712(chainID, accounts.HexToAddress(permitAddress), "Permit2")
	hashed = eip.HashTypedData(hashed)

	err = pp.Sign(s.LocalHardhatMainnetUser.Account, chainID, accounts.HexToAddress(permitAddress), "Permit2")
	s.Require().Nil(err)

	verified, err := s.LocalHardhatMainnetUser.VerifySignature(s.LocalHardhatMainnetUser.Address(), hashed.Bytes(), pp.Signature)
	s.Require().Nil(err)
	s.Require().True(verified)

	// this is why solidity and its idiotic js ecosystem is fucking stupid
	jsSig := "1a622a5fb555e46f58b11ace6176bfc6d1f8ac4be3711612e5f89027de9aae96490d65fc3dce716c08cef58f1d78856fa0a50d13512cd207206d7aca11017ed11b"

	jsSigBytes := pp.Signature
	jsSigBytes[64] += 27
	s.Equal(jsSig, common.Bytes2Hex(jsSigBytes))
}

func (s *Web3ClientTestSuite) TestPermit2Approve() {
	err := s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, 17461070)
	s.Require().Nil(err)
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)

	tx, err := uni.ApproveSpender(ctx, WETH9ContractAddress, Permit2SmartContractAddress, EtherMultiple(10000))
	s.Assert().NoError(err)
	s.Assert().NotNil(tx)

	expiration, _ := new(big.Int).SetString("1678544408", 10)
	sigDeadline, _ := new(big.Int).SetString("1675954208", 10)

	pp := Permit2PermitParams{
		PermitSingle: PermitSingle{
			PermitDetails: PermitDetails{
				Token:      accounts.HexToAddress(WETH9ContractAddress),
				Amount:     EtherMultiple(1),
				Expiration: expiration,
				Nonce:      new(big.Int).SetUint64(0),
			},
			Spender:     accounts.HexToAddress(UniswapUniversalRouterAddressNew),
			SigDeadline: sigDeadline,
		},
		Signature: nil,
	}

	err = pp.Sign(s.LocalHardhatMainnetUser.Account, chainID, accounts.HexToAddress(Permit2SmartContractAddress), "Permit2")
	s.Assert().NoError(err)
	s.Assert().NotNil(pp.Signature)

	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: Permit2SmartContractAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractFile:      "",
		ContractABI:       Permit2AbiDecoder,
		MethodName:        permit0,
		Params:            []interface{}{s.LocalHardhatMainnetUser.Address(), pp.PermitSingle, pp.Signature},
	}

	signedTx, err := s.LocalHardhatMainnetUser.CallFunctionWithArgs(ctx, &params)
	s.Require().Nil(err)
	s.Require().NotNil(signedTx)
}

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
	encCmd, err := ur.EncodeCommands(ctx, nil)
	s.Require().NoError(err)
	s.Require().NotNil(encCmd)
	s.Require().NotNil(encCmd.Commands)
	subCmd := UniversalRouterExecSubCmd{}
	for i, byteVal := range encCmd.Commands {
		err = subCmd.DecodeCommand(byteVal, encCmd.Inputs[i], nil)
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
				To:     accounts.HexToAddress(UniswapUniversalRouterAddressNew),
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
	encCmd, err := ur.EncodeCommands(ctx, nil)
	s.Require().NoError(err)
	s.Require().NotNil(encCmd)
	s.Require().NotNil(encCmd.Commands)
	subCmd := UniversalRouterExecSubCmd{}
	for i, byteVal := range encCmd.Commands {
		err = subCmd.DecodeCommand(byteVal, encCmd.Inputs[i], nil)
		s.Require().NoError(err)
		s.Assert().Equal(true, subCmd.CanRevert)
		s.Assert().Equal(Permit2TransferFromBatch, subCmd.Command)
		decodedInputs := subCmd.DecodedInputs.(Permit2PermitTransferFromBatchParams)
		s.Assert().Equal(permit2TransferFromBatch.Details[0].From, decodedInputs.Details[0].From)
	}
}
