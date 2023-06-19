package web3_client

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (u *UniswapClient) ExecFrontRunTrade(tf TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, &tf.FrontRunTrade)
}

func (u *UniswapClient) ExecFrontRunTradeStep(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	if u.DebugPrint {
		fmt.Println("executing front run trade")
	}
	return u.ExecSwap(tf.InitialPair, &tf.FrontRunTrade)
}

func (u *UniswapClient) ExecFrontRunTradeStepTokenTransfer(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	if u.DebugPrint {
		fmt.Println("executing front run token transfer")
	}
	_, _ = u.FrontRunTradeGetAmountsOut(tf)
	startEthBal, err := u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting pre trade eth balance")
		return nil, err
	}

	bal, _ := u.Web3Client.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountOutAddr.String(), u.Web3Client.PublicKey())
	tf.FrontRunTrade.PreTradeTokenBalance = bal
	fmt.Println("pre trade amount out token balance", bal.String())

	tf.FrontRunTrade.PreTradeEthBalance = startEthBal
	err = u.ExecTradeV2SwapFromTokenToToken(ctx, &tf.FrontRunTrade)
	if err != nil {
		return nil, err
	}

	endEthBal, err := u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting post trade eth balance")
		return nil, err
	}

	bal, _ = u.Web3Client.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountOutAddr.String(), u.Web3Client.PublicKey())
	tf.FrontRunTrade.PostTradeTokenBalance = bal
	fmt.Println("post trade amount out token balance", bal.String())
	tf.FrontRunTrade.PostTradeTokenBalance = bal
	tf.FrontRunTrade.DiffTradeTokenBalance = new(big.Int).Sub(tf.FrontRunTrade.PostTradeTokenBalance, tf.FrontRunTrade.PreTradeTokenBalance)
	fmt.Println("diff trade token balance", tf.FrontRunTrade.DiffTradeTokenBalance.String())
	if tf.FrontRunTrade.AmountOut.String() != tf.FrontRunTrade.DiffTradeTokenBalance.String() {
		return nil, errors.New("balance change does not match prediction")
	}

	tf.FrontRunTrade.PostTradeEthBalance = endEthBal
	tf.FrontRunTrade.DiffTradeEthBalance = new(big.Int).Sub(endEthBal, startEthBal)
	return nil, nil
	//return u.ExecFrontRunTradeStep(tf)
}

func (u *UniswapClient) FrontRunTradeGetAmountsOut(tf *TradeExecutionFlow) ([]*big.Int, error) {
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
		fmt.Println("front run trade expected amount in", tf.FrontRunTrade.AmountIn.String(), "amount out", tf.FrontRunTrade.AmountOut.String())
		fmt.Println("front run trade simulated amount in", amountsOutFirstPair[0].String(), "amount out", amountsOutFirstPair[1].String())
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
