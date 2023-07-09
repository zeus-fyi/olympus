package web3_client

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (s *Web3ClientTestSuite) TestUniversalRouterV2() {
	uni := InitUniswapClient(ctx, s.ProxyHostedHardhatMainnetUser)
	uni.Web3Client.IsAnvilNode = true
	uni.DebugPrint = true
	uni.PrintLocal = true
	uni.PrintDetails = true
	uni.Web3Client.AddSessionLockHeader(uuid.New().String())
	err := uni.Web3Client.HardHatResetNetwork(ctx, 16591736)
	s.Require().Nil(err)

	hashStr := "0x889b34a27b730dd664cd71579b4310522c3b495fb34f17f08d1131c0cec651fa"
	tx, _, err := s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)
	s.Require().NotNil(tx)
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapUniversalRouterNew)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)
	subCmds, err := NewDecodedUniversalRouterExecCmdFromMap(args)
	s.Require().Nil(err)
	s.Require().NotEmpty(subCmds)

	amountIn := ""
	for _, cmd := range subCmds.Commands {
		fmt.Println(cmd.Command)
		if cmd.Command == WrapETH {
			dec := cmd.DecodedInputs.(WrapETHParams)
			fmt.Println("recipient", dec.Recipient.String())
			fmt.Println("amountMin", dec.AmountMin.String())
			amountIn = dec.AmountMin.String()
		}
		if cmd.Command == V2SwapExactOut {
			dec := cmd.DecodedInputs.(V2SwapExactOutParams)
			for _, pa := range dec.Path {
				fmt.Println(pa.String())
			}
			fmt.Println("to", dec.To.String())
			fmt.Println("payerIsSender", dec.PayerIsSender)
			fmt.Println("amountOut", dec.AmountOut.String())
			fmt.Println("amountInMax", dec.AmountInMax.String())
			cmd.CanRevert = false
		}
		if cmd.Command == UnwrapWETH {
			dec := cmd.DecodedInputs.(UnwrapWETHParams)
			fmt.Println("recipient", dec.Recipient.String())
			fmt.Println("amountMin", dec.AmountMin.String())
			cmd.CanRevert = false
		}
	}
	pl, _ := new(big.Int).SetString(amountIn, 10)
	wethParams := WrapETHParams{
		Recipient: s.ProxyHostedHardhatMainnetUser.Address(),
		AmountMin: pl,
	}
	payable := &web3_actions.SendEtherPayload{
		TransferArgs: web3_actions.TransferArgs{
			Amount:    pl,
			ToAddress: wethParams.Recipient,
		},
		GasPriceLimits: web3_actions.GasPriceLimits{},
	}
	subCmds.Payable = payable
	signedTx, err := uni.ExecUniswapUniversalRouterCmd(subCmds)
	s.Require().Nil(err)
	s.Require().NotNil(signedTx)
}

// encodes a single exactInput USDC->ETH swap with permit
func (s *Web3ClientTestSuite) TestV2EthToUsdcSwapWithPermit() {
	expiration, _ := new(big.Int).SetString("3000000000000", 10)
	sigDeadline, _ := new(big.Int).SetString("3000000000000", 10)
	amount := new(big.Int).SetUint64(1000000000)
	usdcAddr := accounts.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	wethAddr := accounts.HexToAddress(WETH9ContractAddress)
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

	verified, err := s.LocalHardhatMainnetUser.VerifySignature(s.LocalHardhatMainnetUser.Address(), hashed.Bytes(), pp.Signature)
	s.Require().Nil(err)
	s.Require().True(verified)

	jsSig := "1a622a5fb555e46f58b11ace6176bfc6d1f8ac4be3711612e5f89027de9aae96490d65fc3dce716c08cef58f1d78856fa0a50d13512cd207206d7aca11017ed100"
	s.Equal(jsSig, common.Bytes2Hex(pp.Signature))

	amountOut, _ := new(big.Int).SetString("780012290817539937", 10)
	v2Trade := V2SwapExactInParams{
		AmountIn:     amount,
		AmountOutMin: amountOut,
		Path: []accounts.Address{
			usdcAddr, wethAddr,
		},
		To:            accounts.HexToAddress("0x0000000000000000000000000000000000000002"),
		PayerIsSender: true,
	}
	scTrade := UniversalRouterExecSubCmd{
		Command:       V2SwapExactIn,
		CanRevert:     false,
		DecodedInputs: v2Trade,
	}
	scPermit := UniversalRouterExecSubCmd{
		Command:       Permit2Permit,
		CanRevert:     false,
		DecodedInputs: pp,
	}
	// export const TEST_RECIPIENT_ADDRESS = '0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
	uwEthParams := UnwrapWETHParams{
		Recipient: accounts.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		AmountMin: amountOut,
	}
	unwrapEth := UniversalRouterExecSubCmd{
		Command:       UnwrapWETH,
		CanRevert:     false,
		DecodedInputs: uwEthParams,
	}
	ep := UniversalRouterExecCmd{
		Commands: []UniversalRouterExecSubCmd{
			scPermit, scTrade, unwrapEth,
		},
		Deadline: nil,
		Payable:  nil,
	}

	data, err := ep.EncodeCommands(ctx)
	s.Require().Nil(err)
	s.Require().NotNil(data)

	scInfo := GetUniswapUniversalRouterAbiPayload(data)
	signedTx, err := s.LocalHardhatMainnetUser.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	s.Require().Nil(err)
	s.Require().NotNil(signedTx)
}
