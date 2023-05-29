package web3_client

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (u *UniswapV2Client) ExecFrontRunTrade(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, &tf.FrontRunTrade)
}

func (u *UniswapV2Client) ExecSandwichTrade(tf TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, &tf.SandwichTrade)
}

func (u *UniswapV2Client) ExecFrontRunTradeStepTokenTransfer(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	if u.DebugPrint {
		fmt.Println("executing front run trade")
	}
	_, _ = u.FrontRunTradeGetAmountsOut(tf)
	ethBal, err := u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting pre trade eth balance")
		return nil, err
	}
	tf.FrontRunTrade.PreTradeEthBalance = ethBal
	err = u.RouterApproveAndSend(ctx, &tf.FrontRunTrade, tf.InitialPair.PairContractAddr)
	if err != nil {
		return nil, err
	}
	ethBal, err = u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting post trade eth balance")
		return nil, err
	}
	tf.FrontRunTrade.PostTradeEthBalance = ethBal
	return u.ExecFrontRunTradeStep(tf)
}

func (u *UniswapV2Client) ExecUserTradeStep(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	if u.DebugPrint {
		fmt.Println("executing user trade")
	}
	_, _ = u.UserTradeGetAmountsOut(tf)
	sender := types.LatestSignerForChainID(tf.Tx.ChainId())
	from, err := sender.Sender(tf.Tx)
	if err != nil {
		log.Err(err).Msg("error getting sender")
		return nil, err
	}
	ethBal, err := u.Web3Client.GetBalance(ctx, from.String(), nil)
	if err != nil {
		log.Err(err).Msg("error getting pre trade eth balance")
		return nil, err
	}
	tf.UserTrade.PreTradeEthBalance = ethBal
	scInfo, err := u.ExecTradeByMethod(tf)
	if err != nil {
		return nil, err
	}
	ethBal, err = u.Web3Client.GetBalance(ctx, from.String(), nil)
	if err != nil {
		log.Err(err).Msg("error getting ending eth balance")
		return nil, err
	}
	tf.UserTrade.PostTradeEthBalance = ethBal
	tf.UserTrade.AddTxHash(accounts.Hash(tf.Tx.Hash()))
	return scInfo, err
}

func (u *UniswapV2Client) ExecFrontRunTradeStep(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	if u.DebugPrint {
		fmt.Println("executing front run trade")
	}
	return u.ExecSwap(tf.InitialPair, &tf.FrontRunTrade)
}

func (u *UniswapV2Client) ExecSandwichTradeStepTokenTransfer(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	_, _ = u.SandwichTradeGetAmountsOut(tf)
	ethBal, err := u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting pre trade eth balance")
		return nil, err
	}
	tf.SandwichTrade.PreTradeEthBalance = ethBal
	err = u.RouterApproveAndSend(ctx, &tf.SandwichTrade, tf.InitialPair.PairContractAddr)
	if err != nil {
		return nil, err
	}
	ethBal, err = u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting post trade eth balance")
		return nil, err
	}
	tf.SandwichTrade.PostTradeEthBalance = ethBal
	return u.ExecSandwichTradeStep(tf)
}

func (u *UniswapV2Client) ExecSandwichTradeStep(tf *TradeExecutionFlowInBigInt) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, &tf.SandwichTrade)
}

func (u *UniswapV2Client) FrontRunTradeGetAmountsOut(tf *TradeExecutionFlowInBigInt) ([]*big.Int, error) {
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
		if u.DebugPrint {
			diff := new(big.Int).Sub(amountsOutFirstPair[1], tf.FrontRunTrade.AmountOut)
			fmt.Println("front run trade actual - expected ", diff.String())
		}
		return amountsOutFirstPair, errors.New("amount out not equal to expected")
	}
	tf.FrontRunTrade.SimulatedAmountOut = amountsOutFirstPair[1]
	return amountsOutFirstPair, err
}

func (u *UniswapV2Client) UserTradeGetAmountsOut(tf *TradeExecutionFlowInBigInt) ([]*big.Int, error) {
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
		if u.DebugPrint {
			diff := new(big.Int).Sub(amountsOutFirstPair[1], tf.UserTrade.AmountOut)
			fmt.Println("user trade actual - expected ", diff.String())
		}
		return amountsOutFirstPair, errors.New("amount out not equal to expected")
	}
	tf.UserTrade.SimulatedAmountOut = amountsOutFirstPair[1]
	return amountsOutFirstPair, err
}

func (u *UniswapV2Client) SandwichTradeGetAmountsOut(tf *TradeExecutionFlowInBigInt) ([]*big.Int, error) {
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
	tf.SandwichTrade.SimulatedAmountOut = amountsOutFirstPair[1]
	return amountsOutFirstPair, err
}
