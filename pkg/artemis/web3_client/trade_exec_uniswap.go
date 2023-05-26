package web3_client

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	web3_types "github.com/zeus-fyi/gochain/web3/types"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

// TODO, low level swap implementation
// Swap (index_topic_1 address sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, index_topic_2 address to)
// Sync (uint112 reserve0, uint112 reserve1)

func (u *UniswapV2Client) ExecFrontRunTrade(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	//tokenNum := tf.InitialPair.GetTokenNumber(tf.FrontRunTrade.AmountInAddr)

	//amount0In := tf.FrontRunTrade.AmountIn.String()
	//amount1In := "0"
	//amount0Out := "0"
	//amount1Out := tf.FrontRunTrade.AmountOut.String()
	//
	//
	//if tokenNum == 1 {
	//	amount0In = "0"
	//	amount1In = tf.FrontRunTrade.AmountIn.String()
	//	amount0Out = tf.FrontRunTrade.AmountOut.String()
	//	amount1Out = "0"
	//}
	scInfo, _, err := LoadSwapAbiPayload(tf.InitialPair.PairContractAddr)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	fmt.Println("00", tf.FrontRunTrade.AmountIn.String())
	fmt.Println("00", tf.FrontRunTrade.AmountInAddr.String())

	fmt.Println("out", tf.FrontRunTrade.AmountOut.String())
	fmt.Println("out", tf.FrontRunTrade.AmountOutAddr.String())

	b, err := u.Web3Client.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountInAddr.String(), u.Web3Client.PublicKey())
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	fmt.Println("balance", b.String())
	ethBal, err := u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	fmt.Println("eth balance", ethBal.String())
	scInfo.Params = []interface{}{tf.FrontRunTrade.AmountOut.String(), tf.FrontRunTrade.AmountIn.String(), u.Web3Client.Address(), []byte{}}
	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	return &scInfo, err
}

/*
Given an output asset amount and an array of token addresses, calculates all preceding minimum
input token amounts by calling getReserves for each pair of token addresses in the path in turn,
and using these to call getAmountIn.
*/
func (u *UniswapV2Client) GetAmountsOut(amountIn *big.Int, pathSlice []string) ([]interface{}, error) {
	pathString := "[" + strings.Join(pathSlice, ",") + "]"
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.MevSmartContractTxMap.Abi,
		MethodName:        getAmountsOut,
		Params:            []interface{}{amountIn, pathString},
	}
	amountsOut, err := u.Web3Client.GetContractConst(ctx, scInfo)
	if err != nil {
		return nil, err
	}
	return amountsOut, err
}

func (u *UniswapV2Client) GetAmountsOutFrontRunTrade(tf TradeExecutionFlowInBigInt) ([]*big.Int, error) {
	pathSlice := []string{tf.FrontRunTrade.AmountInAddr.String(), tf.FrontRunTrade.AmountOutAddr.String()}
	amountsOut, err := u.GetAmountsOut(tf.FrontRunTrade.AmountIn, pathSlice)
	amountsOutFirstPair := ConvertAmountsToBigIntSlice(amountsOut)
	if len(amountsOutFirstPair) != 2 {
		return nil, errors.New("amounts out not equal to expected")
	}
	if tf.FrontRunTrade.AmountIn.String() != amountsOutFirstPair[0].String() {
		return nil, errors.New("amount in not equal to expected")
	}
	if tf.FrontRunTrade.AmountOut.String() != amountsOutFirstPair[1].String() {
		return nil, errors.New("amount out not equal to expected")
	}
	return amountsOutFirstPair, err
}

func ConvertAmountsToBigIntSlice(amounts []interface{}) []*big.Int {
	var amountsBigInt []*big.Int
	for _, amount := range amounts {
		pair := amount.([]*big.Int)
		for _, p := range pair {
			amountsBigInt = append(amountsBigInt, p)
		}
	}
	return amountsBigInt
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
	case swap:
		return u.ExecFrontRunTrade(tf)
	default:
	}
	return nil, errors.New("invalid trade method")
}

func (u *UniswapV2Client) GetAmounts(to TradeOutcome, method string) ([]interface{}, error) {
	switch method {
	case getAmountsOut:
		pathSlice := []string{to.AmountInAddr.String(), to.AmountOutAddr.String()}
		return u.GetAmountsOut(to.AmountIn, pathSlice)
	case getAmountsIn:
		pathSlice := []string{to.AmountOutAddr.String(), to.AmountInAddr.String()}
		return u.GetAmountsIn(to.AmountOut, pathSlice)
	}
	return nil, errors.New("invalid method")
}

/*
Given an output asset amount and an array of token addresses, calculates all preceding minimum
input token amounts by calling getReserves for each pair of token addresses in the path in turn,
and using these to call getAmountIn.
*/
func (u *UniswapV2Client) GetAmountsIn(amountOut *big.Int, pathSlice []string) ([]interface{}, error) {
	pathString := "[" + strings.Join(pathSlice, ",") + "]"
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: u.MevSmartContractTxMap.SmartContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.MevSmartContractTxMap.Abi,
		MethodName:        getAmountsIn,
		Params:            []interface{}{amountOut, pathString},
	}
	amountsIn, err := u.Web3Client.GetContractConst(ctx, scInfo)
	if err != nil {
		return nil, err
	}
	return amountsIn, err
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
