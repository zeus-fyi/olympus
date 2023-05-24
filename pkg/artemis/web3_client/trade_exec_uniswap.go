package web3_client

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

func (w *Web3Client) ExecSwapTrade(tf TradeExecutionFlowInBigInt) (*big.Int, *big.Int) {
	sellAmount := big.NewInt(0)
	maxProfit := big.NewInt(0)
	paramsTx, _, err := LoadSwapAbiPayload()
	if err != nil {
		return sellAmount, maxProfit
	}
	// Pair address in contract
	pairContract := ""
	paramsTx.Params = []interface{}{pairContract, tf.FrontRunTrade.AmountIn, tf.FrontRunTrade.AmountOut}
	return sellAmount, maxProfit
}

func (u *UniswapV2Client) SwapExactTokensForETHParams(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	trade := tf.Trade
	params := *trade.JSONSwapExactTokensForETHParams
	//	pairAddr := u.GetPairContractFromFactory(ctx, params.Path[0].String(), params.Path[1].String())

	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.MevSmartContractTxMap.Abi,
		MethodName:        swapExactTokensForETH,
		Params:            []interface{}{params.AmountIn, params.AmountOutMin, params.Path, params.To, params.Deadline},
	}
	err := u.Web3Client.ImpersonateAccount(ctx, tf.Tx.From.String())
	if err != nil {
		return nil, err
	}
	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, scInfo)
	if err != nil {
		return nil, err
	}
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}
	return scInfo, nil
}
