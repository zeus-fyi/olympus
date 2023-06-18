package web3_client

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

/*
trx_hash_02 = HexStr("0x3247555a5dbc877ade17c4b49362bc981af5fb5064e0b3cbd91411e085fe3093")
expected_function_names_02 = ("V3_SWAP_EXACT_IN", "UNWRAP_WETH")

trx_hash_03 = HexStr("0x889b34a27b730dd664cd71579b4310522c3b495fb34f17f08d1131c0cec651fa")
expected_function_names_03 = ("WRAP_ETH", "V2_SWAP_EXACT_OUT", "UNWRAP_WETH")

trx_hash_04 = HexStr("0xf99ac4237df313794747759550db919b37d7c8a67d4a7e12be8f5cbaacd51376")
expected_function_names_04 = ("WRAP_ETH", "V2_SWAP_EXACT_OUT", "V3_SWAP_EXACT_OUT", "UNWRAP_WETH")
*/

func (s *Web3ClientTestSuite) TestUniversalRouterV3ExactIn() {
	uni := InitUniswapClient(ctx, s.MainnetWeb3User)
	uni.PrintOn = true
	uni.PrintLocal = true
	uni.Path = filepaths.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "./trade_analysis",
		FnIn:        "",
		FnOut:       "",
		Env:         "",
	}
	hashStr := "0x3247555a5dbc877ade17c4b49362bc981af5fb5064e0b3cbd91411e085fe3093"
	tx, _, err := s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)

	s.Require().NotNil(tx)
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapUniversalRouter)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)
	subCmds, err := NewDecodedUniversalRouterExecCmdFromMap(args)
	s.Require().Nil(err)
	s.Require().NotEmpty(subCmds)

	for _, val := range subCmds.Commands {
		fmt.Println(val.Command)
	}
}

func (s *Web3ClientTestSuite) TestUniversalRouterV3ExactOut() {
	uni := InitUniswapClient(ctx, s.MainnetWeb3User)
	uni.PrintOn = true
	uni.PrintLocal = true
	uni.Path = filepaths.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "./trade_analysis",
		FnIn:        "",
		FnOut:       "",
		Env:         "",
	}
	hashStr := "0xf99ac4237df313794747759550db919b37d7c8a67d4a7e12be8f5cbaacd51376"
	tx, _, err := s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)

	s.Require().NotNil(tx)
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapUniversalRouter)
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
				fmt.Println(pa)
			}
			fmt.Println("to", dec.To.String())
			fmt.Println("payerIsSender", dec.PayerIsSender)
			fmt.Println("amountOut", dec.AmountOut.String())
			fmt.Println("amountInMax", dec.AmountInMax.String())
			cmd.CanRevert = false
		}
		if cmd.Command == V3SwapExactOut {
			// TODO fix this broken encoder
			dec := cmd.DecodedInputs.(V3SwapExactOutParams)
			for _, pa := range dec.Path.Path {
				fmt.Println(pa)
			}
			fmt.Println("to", dec.To.String())
			fmt.Println("payerIsUser", dec.PayerIsUser)
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
		Recipient: s.LocalHardhatMainnetUser.Address(),
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
	data, err := subCmds.EncodeCommands(ctx)
	s.Require().Nil(err)
	s.Require().NotNil(data)

	scInfo := GetUniswapUniversalRouterAbiPayload(data)
	signedTx, err := s.LocalHardhatMainnetUser.CallFunctionWithArgs(ctx, &scInfo)
	s.Require().Nil(err)
	s.Require().NotNil(signedTx)
}
