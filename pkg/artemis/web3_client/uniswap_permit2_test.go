package web3_client

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

var (
	chainID = big.NewInt(1)
	urAddr  = accounts.HexToAddress(UniswapUniversalRouterAddress)
)

func (s *Web3ClientTestSuite) TestPermit2Approve() {
	node := "https://virulent-alien-cloud.quiknode.pro/fa84e631e9545d76b9e1b1c5db6607fedf3cb654"
	err := s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, node, 17461070)
	s.Require().Nil(err)
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)

	tx, err := uni.ApproveSpender(ctx, WETH9ContractAddress, Permit2SmartContractAddress, EtherMultiple(10000))
	s.Assert().NoError(err)
	s.Assert().NotNil(tx)

	//
	// 0x889b34a27b730dd664cd71579b4310522c3b495fb34f17f08d1131c0cec651fa
	// 16591736
	// V2_SWAP_EXACT_OUT
	// 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2 -> 0xDadb4aE5B5D3099Dd1f586f990B845F2404A1c4c
	hashStr := "0x52e63b75f41a352ad9182f9e0f923c8557064c3b1047d1778c1ea5b11b979dd9"
	tx, _, err = s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)
	s.Require().NotNil(tx)
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapUniversalRouter)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)
	subCmds, err := NewDecodedUniversalRouterExecCmdFromMap(args)
	s.Require().Nil(err)
	s.Require().NotEmpty(subCmds)
	expiration, _ := new(big.Int).SetString("1678544408", 10)
	sigDeadline, _ := new(big.Int).SetString("1675954208", 10)

	pp := Permit2PermitParams{
		PermitSingle: PermitSingle{
			PermitDetails: PermitDetails{
				TokenPermissions: TokenPermissions{
					Token:  accounts.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
					Amount: EtherMultiple(1),
				},
				Expiration: expiration,
				Nonce:      new(big.Int).SetUint64(0),
			},
			Spender:     accounts.HexToAddress(UniswapUniversalRouterAddress),
			SigDeadline: sigDeadline,
		},
		Signature: nil,
	}

	var sig []byte
	for _, sb := range subCmds.Commands {
		if sb.Command == "PERMIT2_PERMIT" {
			tmp := sb.DecodedInputs.(Permit2PermitParams)

			pp.Token = tmp.Token
			pp.Amount = tmp.Amount
			pp.Expiration = tmp.Expiration
			pp.Nonce = tmp.Nonce
			pp.Spender = tmp.Spender
			pp.SigDeadline = tmp.SigDeadline
			pp.Signature = tmp.Signature
			fmt.Println("token", tmp.Token.Hex())
			fmt.Println("amount", tmp.Expiration.String())
			fmt.Println("expiration", tmp.Expiration.String())
			fmt.Println("nonce", tmp.Nonce.String())
			fmt.Println("spender", tmp.Spender.Hex())
			fmt.Println("sigDeadline", tmp.SigDeadline.String())

			eip := NewEIP712ForPermit2(chainID, urAddr)

			hashedData := hashPermitTransferFrom(pp.TokenPermissions, s.LocalHardhatMainnetUser.Address(), pp.Nonce, pp.Expiration)
			fmt.Println("hashedData", hashedData.Hex())
			result := eip.HashTypedData(hashedData)

			verified, err := s.LocalHardhatMainnetUser.VerifySignature(accounts.HexToAddress("0x2ABa70F8bb7588AAA5CC180BC1f0c6a0b28A4910"), hashedData[:], sig)
			s.Require().Nil(err)
			s.Assert().True(verified)

			verified, err = s.LocalHardhatMainnetUser.VerifySignature(accounts.HexToAddress("0x2ABa70F8bb7588AAA5CC180BC1f0c6a0b28A4910"), result[:], sig)
			s.Require().Nil(err)
			s.Assert().True(verified)
		}

	}

	eip := NewEIP712ForPermit2(chainID, urAddr)

	hashedData := hashPermitTransferFrom(pp.TokenPermissions, s.LocalHardhatMainnetUser.Address(), pp.Nonce, pp.Expiration)
	fmt.Println("hashedData", hashedData.Hex())
	result := eip.HashTypedData(hashedData)

	signed, err := s.LocalHardhatMainnetUser.Sign(result[:])
	s.Require().Nil(err)
	verified, err := s.LocalHardhatMainnetUser.VerifySignature(s.LocalHardhatMainnetUser.Address(), result[:], signed)
	s.Require().Nil(err)
	s.Require().True(verified)

	/*
		    function hash(ISignatureTransfer.PermitTransferFrom memory permit) internal view returns (bytes32) {
		        bytes32 tokenPermissionsHash = _hashTokenPermissions(permit.permitted);
		        return keccak256(
		            abi.encode(_PERMIT_TRANSFER_FROM_TYPEHASH, tokenPermissionsHash, msg.sender, permit.nonce, permit.deadline)
		        );
		    }

			token 0xdAC17F958D2ee523a2206206994597C13D831ec7
			amount 1678544408
			expiration 1678544408
			nonce 0
			spender 0xEf1c6E67703c7BD7107eed8303Fbe6EC2554BF6B
			sigDeadline 1675954208
	*/

	//params := web3_actions.SendContractTxPayload{
	//	SmartContractAddr: Permit2SmartContractAddress,
	//	SendEtherPayload:  web3_actions.SendEtherPayload{},
	//	ContractFile:      "",
	//	ContractABI:       Permit2AbiDecoder,
	//	MethodName:        "permitTransferFrom",
	//	Params:            []interface{}{permitTransferFromParams.Permit, permitTransferFromParams.TransferDetails, permitTransferFromParams.Owner, permitTransferFromParams.Signature},
	//}
	//
	//tx, err = s.LocalHardhatMainnetUser.SignAndSendSmartContractTxPayload(ctx, params)
	//s.Assert().NoError(err)
	//s.Assert().NotNil(tx)

	// todo prove permit2 transfer works natively, then test via UR
}

/*
  function DOMAIN_SEPARATOR() public view override returns (bytes32) {
        return
            keccak256(
                abi.encode(
                    // keccak256('EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)')
                    0x8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f,
                    nameHash,
                    versionHash,
                    ChainId.get(),
                    address(this)
                )
            );
    }
*/
/*
			   PermitTransferFrom memory permit,
			   SignatureTransferDetails calldata transferDetails,
			   address owner,
			   bytes calldata signature


			/// @notice The signed permit message for a single token transfer
			struct PermitTransferFrom {
				TokenPermissions permitted;
				// a unique value for every token owner's signature to prevent signature replays
				uint256 nonce;
				// deadline on the permit signature
				uint256 deadline;
			}

			/// @notice The token and amount details for a transfer signed in the permit transfer signature
			struct TokenPermissions {
				// ERC20 token address
				address token;
				// the maximum amount that can be spent
				uint256 amount;
			}

	    /// @notice Specifies the recipient address and amount for batched transfers.
	    /// @dev Recipients and amounts correspond to the index of the signed token permissions array.
	    /// @dev Reverts if the requested amount is greater than the permitted signed amount.
		    struct SignatureTransferDetails {
		        // recipient address
		        address to;
		        // spender requested amount
		        uint256 requestedAmount;
		    }
*/

func (s *Web3ClientTestSuite) TestPermit2PermitBatchEncode() {
	addr1 := accounts.HexToAddress(LidoSEthAddr)
	addr2 := accounts.HexToAddress(WETH9ContractAddress)
	permit2Batch := Permit2PermitBatchParams{
		PermitBatch: PermitBatch{
			Details: []PermitDetails{{
				TokenPermissions: TokenPermissions{
					Token:  addr1,
					Amount: new(big.Int).SetUint64(1000000000000000000),
				},
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
