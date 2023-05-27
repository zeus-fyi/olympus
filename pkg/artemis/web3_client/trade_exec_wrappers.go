package web3_client

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

func (u *UniswapV2Client) ExecFrontRunTradeStepTokenTransfer(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	if u.DebugPrint {
		fmt.Println("executing front run trade")
	}
	err := u.RouterApproveAndSend(ctx, &tf.FrontRunTrade, tf.InitialPair.PairContractAddr)
	if err != nil {
		return nil, err
	}
	return u.ExecFrontRunTradeStep(tf)
}

func (u *UniswapV2Client) ExecFrontRunTradeStep(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	if u.DebugPrint {
		fmt.Println("executing front run trade")
	}
	return u.ExecSwap(tf.InitialPair, &tf.FrontRunTrade)
}

func (u *UniswapV2Client) ExecSandwichTradeStepTokenTransfer(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	err := u.RouterApproveAndSend(ctx, &tf.SandwichTrade, tf.InitialPair.PairContractAddr)
	if err != nil {
		return nil, err
	}
	return u.ExecSandwichTradeStep(tf)
}

func (u *UniswapV2Client) ExecSandwichTradeStep(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, &tf.SandwichTrade)
}

func (u *UniswapV2Client) FrontRunTradeGetAmountsOut(tf TradeExecutionFlowInBigInt) ([]*big.Int, error) {
	pathSlice := []string{tf.FrontRunTrade.AmountInAddr.String(), tf.FrontRunTrade.AmountOutAddr.String()}
	amountsOut, err := u.GetAmountsOut(tf.FrontRunTrade.AmountIn, pathSlice)
	if err != nil {
		return nil, err
	}
	amountsOutFirstPair := ConvertAmountsToBigIntSlice(amountsOut)
	if len(amountsOutFirstPair) != 2 {
		return nil, errors.New("amounts out not equal to expected")
	}
	if u.DebugPrint {
		fmt.Println("front run trade trade path", pathSlice[0], pathSlice[1])
		fmt.Println("front run trade expected amounts", amountsOutFirstPair[0].String(), amountsOutFirstPair[1].String())
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
	if err != nil {
		return nil, err
	}
	amountsOutFirstPair := ConvertAmountsToBigIntSlice(amountsOut)
	if len(amountsOutFirstPair) != 2 {
		return nil, errors.New("amounts out not equal to expected")
	}
	if u.DebugPrint {
		fmt.Println("user trade trade path", pathSlice[0], pathSlice[1])
		fmt.Println("user trade expected amounts", amountsOutFirstPair[0].String(), amountsOutFirstPair[1].String())
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

func (u *UniswapV2Client) SandwichTradeGetAmountsOut(tf TradeExecutionFlowInBigInt) ([]*big.Int, error) {
	pathSlice := []string{tf.SandwichTrade.AmountInAddr.String(), tf.SandwichTrade.AmountOutAddr.String()}
	amountsOut, err := u.GetAmountsOut(tf.SandwichTrade.AmountIn, pathSlice)
	if err != nil {
		return nil, err
	}
	amountsOutFirstPair := ConvertAmountsToBigIntSlice(amountsOut)
	if len(amountsOutFirstPair) != 2 {
		return nil, errors.New("amounts out not equal to expected")
	}
	if u.DebugPrint {
		fmt.Println("sandwich trade trade path", pathSlice[0], pathSlice[1])
		fmt.Println("sandwich trade expected amounts", amountsOutFirstPair[0].String(), amountsOutFirstPair[1].String())
	}
	if tf.SandwichTrade.AmountIn.String() != amountsOutFirstPair[0].String() {
		log.Warn().Msgf(fmt.Sprintf("amount in not equal to expected amount in %s, actual amount in: %s", tf.UserTrade.AmountIn.String(), amountsOutFirstPair[0].String()))
		return amountsOutFirstPair, errors.New("amount in not equal to expected")
	}
	if tf.SandwichTrade.AmountOut.String() != amountsOutFirstPair[1].String() {
		log.Warn().Msgf(fmt.Sprintf("amount out not equal to expected amount out %s, actual amount out: %s", tf.UserTrade.AmountOut.String(), amountsOutFirstPair[1].String()))
		return amountsOutFirstPair, errors.New("amount out not equal to expected")
	}
	return amountsOutFirstPair, err
}
