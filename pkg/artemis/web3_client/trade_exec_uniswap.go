package web3_client

import (
	"context"
	"errors"
	"strings"

	web3_types "github.com/zeus-fyi/gochain/web3/types"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

// TODO, low level swap implementation
// Swap (index_topic_1 address sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, index_topic_2 address to)
// Sync (uint112 reserve0, uint112 reserve1)

func (u *UniswapV2Client) ExecFrontRunTrade(tf TradeExecutionFlowInBigInt) error {
	amount0In := ""
	amount1In := ""
	amount0Out := ""
	amount1Out := ""
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.MevSmartContractTxMap.Abi,
		MethodName:        swap,
		Params:            []interface{}{amount0In, amount1In, amount0Out, amount1Out},
	}

	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, scInfo)
	if err != nil {
		return err
	}
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return err
	}
	return err
}

func (u *UniswapV2Client) ExecTradeByMethod(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	switch tf.Trade.TradeMethod {
	case swapTokensForExactETH:
		return u.SwapTokensForExactETHParams(tf)
	case swapExactTokensForETH:
		return u.SwapExactTokensForETHParams(tf)
	case swapExactTokensForTokens:
		return u.SwapExactTokensForTokensParams(tf)
	case swapTokensForExactTokens:
		return u.SwapTokensForExactTokensParams(tf)
	case swapExactETHForTokens:
		return u.SwapExactETHForTokensParams(tf)
	case swapETHForExactTokens:
		return u.SwapETHForExactTokensParams(tf)
	default:
	}
	return nil, errors.New("invalid trade method")
}

func (u *UniswapV2Client) SwapExactTokensForETHParams(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	trade := tf.Trade
	params := *trade.JSONSwapExactTokensForETHParams
	pathSlice := make([]string, len(params.Path))
	for i, p := range params.Path {
		pathSlice[i] = p.String()
	}
	pathString := "[" + strings.Join(pathSlice, ",") + "]"
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.MevSmartContractTxMap.Abi,
		MethodName:        swapExactTokensForETH,
		Params:            []interface{}{params.AmountIn, params.AmountOutMin, pathString, params.To, params.Deadline},
	}
	err := u.Web3Client.SendImpersonatedTx(ctx, tf.Tx)
	if err != nil {
		return nil, err
	}
	return scInfo, nil
}

func (u *UniswapV2Client) SwapTokensForExactETHParams(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	trade := tf.Trade
	params := *trade.JSONSwapTokensForExactETHParams
	pathSlice := make([]string, len(params.Path))
	for i, p := range params.Path {
		pathSlice[i] = p.String()
	}
	pathString := "[" + strings.Join(pathSlice, ",") + "]"

	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.MevSmartContractTxMap.Abi,
		MethodName:        swapTokensForExactETH,
		Params:            []interface{}{params.AmountOut, params.AmountInMax, pathString, params.To, params.Deadline},
	}
	err := u.Web3Client.SendImpersonatedTx(ctx, tf.Tx)
	if err != nil {
		return nil, err
	}
	return scInfo, nil
}

func (u *UniswapV2Client) SwapExactTokensForTokensParams(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	trade := tf.Trade
	params := *trade.JSONSwapExactTokensForTokensParams
	pathSlice := make([]string, len(params.Path))
	for i, p := range params.Path {
		pathSlice[i] = p.String()
	}
	pathString := "[" + strings.Join(pathSlice, ",") + "]"

	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.MevSmartContractTxMap.Abi,
		MethodName:        swapExactTokensForTokens,
		Params:            []interface{}{params.AmountIn, params.AmountOutMin, pathString, params.To, params.Deadline},
	}
	err := u.Web3Client.SendImpersonatedTx(ctx, tf.Tx)
	if err != nil {
		return nil, err
	}
	return scInfo, nil
}

func (u *UniswapV2Client) SwapExactETHForTokensParams(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	trade := tf.Trade
	params := *trade.JSONSwapExactETHForTokensParams
	pathSlice := make([]string, len(params.Path))
	for i, p := range params.Path {
		pathSlice[i] = p.String()
	}
	pathString := "[" + strings.Join(pathSlice, ",") + "]"

	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload: web3_actions.SendEtherPayload{
			TransferArgs: web3_actions.TransferArgs{
				Amount:    tf.Tx.Value.ToInt(),
				ToAddress: params.To,
			},
		},
		ContractABI: u.MevSmartContractTxMap.Abi,
		MethodName:  swapExactETHForTokens,
		Params:      []interface{}{params.AmountOutMin, pathString, params.To, params.Deadline},
	}
	err := u.Web3Client.SendImpersonatedTx(ctx, tf.Tx)
	if err != nil {
		return nil, err
	}
	return scInfo, nil
}

func (u *UniswapV2Client) SwapETHForExactTokensParams(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	trade := tf.Trade
	params := *trade.JSONSwapETHForExactTokensParams
	pathSlice := make([]string, len(params.Path))
	for i, p := range params.Path {
		pathSlice[i] = p.String()
	}
	pathString := "[" + strings.Join(pathSlice, ",") + "]"

	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload: web3_actions.SendEtherPayload{
			TransferArgs: web3_actions.TransferArgs{
				Amount:    tf.Tx.Value.ToInt(),
				ToAddress: params.To,
			},
		},
		ContractABI: u.MevSmartContractTxMap.Abi,
		MethodName:  swapETHForExactTokens,
		Params:      []interface{}{params.AmountOut, pathString, params.To, params.Deadline},
	}
	err := u.Web3Client.SendImpersonatedTx(ctx, tf.Tx)
	if err != nil {
		return nil, err
	}
	return scInfo, nil
}

func (u *UniswapV2Client) SwapTokensForExactTokensParams(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	trade := tf.Trade
	params := *trade.JSONSwapTokensForExactTokensParams
	pathSlice := make([]string, len(params.Path))
	for i, p := range params.Path {
		pathSlice[i] = p.String()
	}
	pathString := "[" + strings.Join(pathSlice, ",") + "]"
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.MevSmartContractTxMap.Abi,
		MethodName:        swapTokensForExactTokens,
		Params:            []interface{}{params.AmountOut, params.AmountInMax, pathString, params.To, params.Deadline},
	}
	err := u.Web3Client.SendImpersonatedTx(ctx, tf.Tx)
	if err != nil {
		return nil, err
	}
	return scInfo, nil
}

func (w *Web3Client) SendImpersonatedTx(ctx context.Context, tx *web3_types.RpcTransaction) error {
	err := w.ImpersonateAccount(ctx, tx.From.String())
	if err != nil {
		return err
	}
	err = w.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}
	err = w.StopImpersonatingAccount(ctx, tx.From.String())
	if err != nil {
		return err
	}
	return nil
}
