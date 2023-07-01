package web3_client

func (u *UniswapClient) ProcessV2Router02Txs() {
	for _, tx := range u.MevSmartContractTxMapV2Router02.Txs {
		switch tx.MethodName {
		case addLiquidity:
			//u.AddLiquidity(tx.Args)
		case addLiquidityETH:
			// payable
			//u.AddLiquidityETH(tx.Args)
			if tx.Tx.Value() == nil {
				continue
			}
		case removeLiquidity:
			//u.RemoveLiquidity(tx.Args)
		case removeLiquidityETH:
			//u.RemoveLiquidityETH(tx.Args)
		case removeLiquidityWithPermit:
			//u.RemoveLiquidityWithPermit(tx.Args)
		case removeLiquidityETHWithPermit:
			//u.RemoveLiquidityETHWithPermit(tx.Args)
		case swapExactTokensForTokens:
			u.SwapExactTokensForTokens(tx, tx.Args)
		case swapTokensForExactTokens:
			u.SwapTokensForExactTokens(tx, tx.Args)
		case swapExactETHForTokens:
			// payable
			if tx.Tx.Value() == nil {
				continue
			}
			u.SwapExactETHForTokens(tx, tx.Args, tx.Tx.Value())
		case swapTokensForExactETH:
			u.SwapTokensForExactETH(tx, tx.Args)
		case swapExactTokensForETH:
			u.SwapExactTokensForETH(tx, tx.Args)
		case swapETHForExactTokens:
			// payable
			if tx.Tx.Value() == nil {
				continue
			}
			u.SwapETHForExactTokens(tx, tx.Args, tx.Tx.Value())
		case swapExactTokensForETHSupportingFeeOnTransferTokensMoniker:
			u.SwapExactTokensForETHSupportingFeeOnTransferTokens(tx, tx.Args)
		case swapExactETHForTokensSupportingFeeOnTransferTokens:
			// payable
			if tx.Tx.Value() == nil {
				continue
			}
			u.SwapExactETHForTokensSupportingFeeOnTransferTokensParams(tx, tx.Args, tx.Tx.Value())
		case swapExactTokensForTokensSupportingFeeOnTransferTokens:
			u.SwapExactTokensForTokensSupportingFeeOnTransferTokens(tx, tx.Args)
		}
	}
}
