package web3_client

import (
	"errors"
	"strings"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (u *UniswapClient) ExecTradeByMethod(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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

func (u *UniswapClient) SwapExactTokensForETHParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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

func (u *UniswapClient) SwapTokensForExactETHParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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

func (u *UniswapClient) SwapExactTokensForTokensParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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

func (u *UniswapClient) SwapExactETHForTokensParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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
				Amount:    tf.Tx.Value(),
				ToAddress: accounts.Address(params.To),
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

func (u *UniswapClient) SwapETHForExactTokensParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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
				Amount:    tf.Tx.Value(),
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

func (u *UniswapClient) SwapTokensForExactTokensParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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
