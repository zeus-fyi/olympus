package web3_client

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (u *UniswapClient) ExecSandwichTrade(tf TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, &tf.SandwichTrade)
}

func (u *UniswapClient) ExecSandwichTradeStep(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, &tf.SandwichTrade)
}

func (u *UniswapClient) ExecSandwichTradeStepTokenTransfer(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	_, _ = u.SandwichTradeGetAmountsOut(tf)
	startEthBal, err := u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting pre trade eth balance")
		return nil, err
	}
	bal, err := u.Web3Client.ReadERC20TokenBalance(ctx, tf.SandwichTrade.AmountOutAddr.String(), u.Web3Client.PublicKey())
	if err != nil {
		log.Err(err).Msg("error getting pre trade amount out token balance")
		return nil, err
	}
	tf.SandwichTrade.PreTradeTokenBalance = bal
	fmt.Println("pre trade amount out token balance", bal.String())

	tf.SandwichTrade.PreTradeEthBalance = startEthBal
	err = u.ExecTradeV2SwapFromTokenToToken(ctx, &tf.SandwichTrade)
	if err != nil {
		return nil, err
	}
	endEthBal, err := u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting post trade eth balance")
		return nil, err
	}

	bal, err = u.Web3Client.ReadERC20TokenBalance(ctx, tf.SandwichTrade.AmountOutAddr.String(), u.Web3Client.PublicKey())
	if err != nil {
		log.Err(err).Msg("error getting post trade amount out token balance")
		return nil, err
	}
	tf.SandwichTrade.PostTradeTokenBalance = bal
	tf.SandwichTrade.DiffTradeTokenBalance = new(big.Int).Sub(tf.SandwichTrade.PostTradeTokenBalance, tf.SandwichTrade.PreTradeTokenBalance)
	fmt.Println("actual amount out", tf.SandwichTrade.DiffTradeTokenBalance.String())
	fmt.Println("expected amount out", tf.SandwichTrade.AmountOut.String())
	drift := new(big.Int).Sub(tf.SandwichTrade.AmountOut, tf.SandwichTrade.DiffTradeTokenBalance)
	if drift.String() != "0" {
		return nil, errors.New("balance change does not match prediction")
	}
	tf.SandwichTrade.PostTradeEthBalance = endEthBal
	tf.SandwichTrade.DiffTradeEthBalance = new(big.Int).Sub(endEthBal, startEthBal)
	return nil, nil
}

func (u *UniswapClient) SandwichTradeGetAmountsOut(tf *TradeExecutionFlow) ([]*big.Int, error) {
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
		fmt.Println("sandwich trade expected amount in", tf.SandwichTrade.AmountIn.String(), "amount out", tf.SandwichTrade.AmountOut.String())
		fmt.Println("sandwich trade simulated amount in", amountsOutFirstPair[0].String(), "amount out", amountsOutFirstPair[1].String())
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
