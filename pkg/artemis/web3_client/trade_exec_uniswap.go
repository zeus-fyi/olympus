package web3_client

import (
	"errors"
	"strings"

	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

func (u *UniswapV2Client) ExecTradeByMethod(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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

func (u *UniswapV2Client) ExecSwap(pair UniswapV2Pair, to *TradeOutcome) (*web3_actions.SendContractTxPayload, error) {
	scInfo, _, err := LoadSwapAbiPayload(pair.PairContractAddr)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	tokenNum := pair.GetTokenNumber(to.AmountInAddr)
	if tokenNum == 0 {
		scInfo.Params = []interface{}{"0", to.AmountOut, u.Web3Client.Address(), []byte{}}
	} else {
		scInfo.Params = []interface{}{to.AmountOut, "0", u.Web3Client.Address(), []byte{}}
	}
	// TODO implement better gas estimation
	scInfo.GasLimit = 3000000
	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	to.AddTxHash(signedTx.Hash())
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	return &scInfo, nil
}

func (u *UniswapV2Client) SwapExactTokensForETHParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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

func (u *UniswapV2Client) SwapTokensForExactETHParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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

func (u *UniswapV2Client) SwapExactTokensForTokensParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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

func (u *UniswapV2Client) SwapExactETHForTokensParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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

func (u *UniswapV2Client) SwapETHForExactTokensParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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

func (u *UniswapV2Client) SwapTokensForExactTokensParams(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
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
