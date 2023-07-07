package artemis_trading_auxiliary

/*
func (s *Web3ClientTestSuite) TestUniversalRouterV2() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ForceDirToTestDirLocation()
	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)

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

	node := "https://virulent-alien-cloud.quiknode.pro/fa84e631e9545d76b9e1b1c5db6607fedf3cb654"
	err = s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, node, 16591736)
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
*/
