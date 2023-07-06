package web3_client

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

func (u *UniswapClient) ExecFrontRunTrade(tf TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(*tf.InitialPair, &tf.FrontRunTrade)
}

func (u *UniswapClient) ExecFrontRunTradeStep(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	if u.DebugPrint {
		fmt.Println("executing front run trade")
	}
	return u.ExecSwap(*tf.InitialPair, &tf.FrontRunTrade)
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
	bal, err := u.Web3Client.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountOutAddr.String(), u.Web3Client.PublicKey())
	if err != nil {
		log.Err(err).Msg("error getting pre trade amount out token balance")
		return nil, err
	}
	tf.FrontRunTrade.PreTradeTokenBalance = bal
	fmt.Println("pre trade amount out token balance", bal.String())
	tf.FrontRunTrade.PreTradeEthBalance = startEthBal

	if tf.InitialPairV3 != nil {
		err = u.ExecTradeV3SwapFromTokenToToken(ctx, &tf.FrontRunTrade)
		if err != nil {
			return nil, err
		}
	} else {
		err = u.ExecTradeV2SwapFromTokenToToken(ctx, &tf.FrontRunTrade)
		if err != nil {
			return nil, err
		}
	}

	endEthBal, err := u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting post trade eth balance")
		return nil, err
	}
	bal, err = u.Web3Client.ReadERC20TokenBalance(ctx, tf.FrontRunTrade.AmountOutAddr.String(), u.Web3Client.PublicKey())
	if err != nil {
		log.Err(err).Msg("error getting post trade amount out token balance")
		return nil, err
	}
	tf.FrontRunTrade.PostTradeTokenBalance = bal
	fmt.Println("post trade amount out token balance", bal.String())
	tf.FrontRunTrade.PostTradeTokenBalance = bal
	tf.FrontRunTrade.DiffTradeTokenBalance = new(big.Int).Sub(tf.FrontRunTrade.PostTradeTokenBalance, tf.FrontRunTrade.PreTradeTokenBalance)
	fmt.Println("diff trade token balance", tf.FrontRunTrade.DiffTradeTokenBalance.String())
	if artemis_eth_units.IsXLessThanY(tf.FrontRunTrade.AmountOut, tf.FrontRunTrade.DiffTradeTokenBalance) {
		log.Info().Msgf("amount out %s is less than the diff trade token balance %s", tf.FrontRunTrade.AmountOut.String(), tf.FrontRunTrade.DiffTradeTokenBalance.String())
		percentDiff := artemis_eth_units.PercentDiff(tf.FrontRunTrade.AmountOut, tf.FrontRunTrade.DiffTradeTokenBalance)
		actualDiff := new(big.Int).Sub(tf.FrontRunTrade.AmountOut, tf.FrontRunTrade.DiffTradeTokenBalance)
		log.Info().Msgf("actual diff %s", actualDiff.String())
		log.Info().Msgf("percent diff %s", percentDiff.String())
		return nil, errors.New("balance change does not match prediction")
	}

	tf.FrontRunTrade.PostTradeEthBalance = endEthBal
	tf.FrontRunTrade.DiffTradeEthBalance = new(big.Int).Sub(endEthBal, startEthBal)
	return nil, nil
}

func (u *UniswapClient) FrontRunTradeGetAmountsOut(tf *TradeExecutionFlow) ([]*big.Int, error) {
	pathSlice := []string{tf.FrontRunTrade.AmountInAddr.String(), tf.FrontRunTrade.AmountOutAddr.String()}
	amountsOutFirstPair, err := u.GetAmountsOut(tf.Tx.To(), tf.FrontRunTrade.AmountIn, pathSlice)
	if err != nil {
		return nil, err
	}
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

	if artemis_eth_units.IsXLessThanY(tf.FrontRunTrade.AmountOut, amountsOutFirstPair[1]) {
		log.Warn().Msgf(fmt.Sprintf("amount out not equal to expected amount out %s, actual amount out: %s", tf.FrontRunTrade.AmountOut.String(), amountsOutFirstPair[1].String()))
		diff := new(big.Int).Sub(amountsOutFirstPair[1], tf.FrontRunTrade.AmountOut)
		if u.DebugPrint {
			fmt.Println("front run trade actual - expected ", diff.String())
		}
		tf.FrontRunTrade.AmountOutDrift = diff
		percentDiff := artemis_eth_units.PercentDiff(tf.FrontRunTrade.AmountOut, amountsOutFirstPair[1])
		log.Info().Msgf("percent diff %s", percentDiff.String())
		return amountsOutFirstPair, errors.New("amount out not equal to expected")
	}

	//if tf.FrontRunTrade.AmountOut.String() != amountsOutFirstPair[1].String() {
	//	log.Warn().Msgf(fmt.Sprintf("amount out not equal to expected amount out %s, actual amount out: %s", tf.FrontRunTrade.AmountOut.String(), amountsOutFirstPair[1].String()))
	//	diff := new(big.Int).Sub(amountsOutFirstPair[1], tf.FrontRunTrade.AmountOut)
	//	if u.DebugPrint {
	//		fmt.Println("front run trade actual - expected ", diff.String())
	//	}
	//	tf.FrontRunTrade.AmountOutDrift = diff
	//	return amountsOutFirstPair, errors.New("amount out not equal to expected")
	//}
	tf.FrontRunTrade.SimulatedAmountOut = amountsOutFirstPair[1]
	return amountsOutFirstPair, err
}
