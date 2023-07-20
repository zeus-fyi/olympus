package web3_client

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
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
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapUniversalRouterNew)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)
	subCmds, err := NewDecodedUniversalRouterExecCmdFromMap(args, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(subCmds)

	for _, cmd := range subCmds.Commands {
		fmt.Println(cmd.Command)

		if cmd.Command == V3SwapExactIn {
			dec := cmd.DecodedInputs.(V3SwapExactInParams)
			fmt.Println("tokenIn", dec.Path.TokenIn.String())
			for _, pa := range dec.Path.Path {
				fmt.Println(pa.Token.String())
				fmt.Println(pa.Fee.String())
			}
			fmt.Println("to", dec.To.String())
			fmt.Println("payerIsUser", dec.PayerIsUser)
			fmt.Println("amountIn", dec.AmountIn.String())
			fmt.Println("amountOutMin", dec.AmountOutMin.String())
			cmd.CanRevert = false
		}
	}
}

func (s *Web3ClientTestSuite) TestUniversalRouterV3ExactOut() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
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

	err = s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, 16606334)
	s.Require().Nil(err)
	mn, args, err := DecodeTxArgData(ctx, tx, uni.MevSmartContractTxMapUniversalRouterNew)
	s.Require().Nil(err)
	s.Require().NotEmpty(mn)
	s.Require().NotEmpty(args)
	subCmds, err := NewDecodedUniversalRouterExecCmdFromMap(args, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(subCmds)

	bal, err := s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, "0xf3dcbc6D72a4E1892f7917b7C43b74131Df8480e", s.LocalHardhatMainnetUser.Address().String())
	fmt.Println("bal", bal.String())
	s.Require().Nil(err)

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
		if cmd.Command == V3SwapExactOut {
			dec := cmd.DecodedInputs.(V3SwapExactOutParams)
			fmt.Println("tokenIn", dec.Path.TokenIn.String())
			for _, pa := range dec.Path.Path {
				fmt.Println(pa.Fee.String())
				fmt.Println(pa.Token.String())
			}

			fmt.Println("to", dec.To.String())
			fmt.Println("payerIsUser", dec.PayerIsUser)
			fmt.Println("amountOut", dec.AmountOut.String())
			fmt.Println("amountInMax", dec.AmountInMax.String())
			cmd.CanRevert = false

			tmp := dec.Path
			expPathBytes := "f3dcbc6d72a4e1892f7917b7c43b74131df8480e000bb8c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
			s.Require().Equal(expPathBytes, common.Bytes2Hex(tmp.Encode()))
		}
		if cmd.Command == UnwrapWETH {
			dec := cmd.DecodedInputs.(UnwrapWETHParams)
			fmt.Println("recipient", dec.Recipient.String())
			fmt.Println("amountMin", dec.AmountMin.String())
			cmd.CanRevert = false
		}
	}
	pl, _ := new(big.Int).SetString(amountIn, 10)
	payable := &web3_actions.SendEtherPayload{
		TransferArgs: web3_actions.TransferArgs{
			Amount:    pl,
			ToAddress: s.LocalHardhatMainnetUser.Address(),
		},
		GasPriceLimits: web3_actions.GasPriceLimits{},
	}
	subCmds.Payable = payable

	tx, err = uni.ExecUniswapUniversalRouterCmd(subCmds)
	s.Assert().Nil(err)
	s.Assert().NotNil(tx)

	bal, err = s.LocalHardhatMainnetUser.ReadERC20TokenBalance(ctx, "0xf3dcbc6D72a4E1892f7917b7C43b74131Df8480e", s.LocalHardhatMainnetUser.Address().String())
	fmt.Println("bal", bal.String())
	s.Assert().Nil(err)

	s.Assert().Equal("8000000000000000002659", bal.String())
}
