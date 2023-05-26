package web3_client

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/rs/zerolog/log"
	web3_types "github.com/zeus-fyi/gochain/web3/types"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

// TODO, low level swap implementation
// Swap (index_topic_1 address sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, index_topic_2 address to)
// Sync (uint112 reserve0, uint112 reserve1)

/*
   the swap function is used to trade from one token to another

   token in = the token address you want to trade out of
   token out = the token address you want as the output of this trade
   amount in = the amount of tokens you are sending in
   amount out Min = the minimum amount of tokens you want out of the trade
   to = the address you want the tokens to be sent to
*/

//   function swap(address _tokenIn, address _tokenOut, uint256 _amountIn, uint256 _amountOutMin, address _to) external

func (u *UniswapV2Client) ExecSwap(pair UniswapV2Pair, to TradeOutcome) (*web3_actions.SendContractTxPayload, error) {
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
	scInfo.GasLimit = 3000000
	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	return &scInfo, nil
}

func (u *UniswapV2Client) ExecFrontRunTrade(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, tf.FrontRunTrade)
}

func (u *UniswapV2Client) ExecSandwichTrade(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, tf.FrontRunTrade)
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

func (u *UniswapV2Client) FrontRunTradeGetAmountsOut(tf TradeExecutionFlowInBigInt) ([]*big.Int, error) {
	pathSlice := []string{tf.FrontRunTrade.AmountInAddr.String(), tf.FrontRunTrade.AmountOutAddr.String()}
	amountsOut, err := u.GetAmountsOut(tf.FrontRunTrade.AmountIn, pathSlice)
	amountsOutFirstPair := ConvertAmountsToBigIntSlice(amountsOut)
	if len(amountsOutFirstPair) != 2 {
		return nil, errors.New("amounts out not equal to expected")
	}
	if tf.FrontRunTrade.AmountIn.String() != amountsOutFirstPair[0].String() {
		log.Warn().Msgf(fmt.Sprintf("amount in not equal to expected amount in %s, actual amount in: %s", tf.FrontRunTrade.AmountIn.String(), amountsOutFirstPair[0].String()))
		return amountsOutFirstPair, errors.New("amount in not equal to expected")
	}
	if tf.FrontRunTrade.AmountOut.String() != amountsOutFirstPair[1].String() {
		log.Warn().Msgf(fmt.Sprintf("amount out not equal to expected amount out %s, actual amount out: %s", tf.FrontRunTrade.AmountOut.String(), amountsOutFirstPair[1].String()))
		return amountsOutFirstPair, errors.New("amount out not equal to expected")
	}
	return amountsOutFirstPair, err
}

func (u *UniswapV2Client) UserTradeGetAmountsOut(tf TradeExecutionFlowInBigInt) ([]*big.Int, error) {
	pathSlice := []string{tf.UserTrade.AmountInAddr.String(), tf.UserTrade.AmountOutAddr.String()}
	amountsOut, err := u.GetAmountsOut(tf.UserTrade.AmountIn, pathSlice)
	amountsOutFirstPair := ConvertAmountsToBigIntSlice(amountsOut)
	if len(amountsOutFirstPair) != 2 {
		return nil, errors.New("amounts out not equal to expected")
	}
	if tf.UserTrade.AmountIn.String() != amountsOutFirstPair[0].String() {
		log.Warn().Msgf(fmt.Sprintf("amount in not equal to expected amount in %s, actual amount in: %s", tf.UserTrade.AmountIn.String(), amountsOutFirstPair[0].String()))
		return amountsOutFirstPair, errors.New("amount in not equal to expected")
	}
	if tf.UserTrade.AmountOut.String() != amountsOutFirstPair[1].String() {
		log.Warn().Msgf(fmt.Sprintf("amount out not equal to expected amount out %s, actual amount out: %s", tf.UserTrade.AmountOut.String(), amountsOutFirstPair[1].String()))
		return amountsOutFirstPair, errors.New("amount out not equal to expected")
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
	case swapFrontRun:
		return u.ExecSwap(tf.InitialPair, tf.FrontRunTrade)
	case swapSandwich:
		return u.ExecSwap(tf.InitialPair, tf.SandwichTrade)
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
	//txHash := common.HexToHash("0x6154c2f46973cecb3f4bc4c1508e3271f3e8dbf7cbcd3c3747b699b8b06b8185")
	//rx, err := w.GetTransactionReceipt(ctx, txHash)
	//if err != nil {
	//	return err
	//}
	//// 36627061988 * 114409
	//fmt.Println(tx.GasPrice.ToInt().String())
	//fmt.Println(rx.GasUsed)
	//fmt.Println(rx.CumulativeGasUsed)
	return nil
}
