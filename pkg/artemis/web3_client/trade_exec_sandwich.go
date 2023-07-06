package web3_client

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

func (u *UniswapClient) ExecSandwichTrade(tf TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(*tf.InitialPair, &tf.SandwichTrade)
}

func (u *UniswapClient) ExecSandwichTradeStep(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(*tf.InitialPair, &tf.SandwichTrade)
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
	if tf.InitialPairV3 != nil {
		err = u.ExecTradeV3SwapFromTokenToToken(ctx, &tf.SandwichTrade)
		if err != nil {
			return nil, err
		}
	} else {
		err = u.ExecTradeV2SwapFromTokenToToken(ctx, &tf.SandwichTrade)
		if err != nil {
			return nil, err
		}
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
	if artemis_eth_units.IsXLessThanY(tf.SandwichTrade.AmountOut, tf.SandwichTrade.DiffTradeTokenBalance) {
		log.Info().Msgf("amount out %s is less than the diff trade token balance %s", tf.SandwichTrade.AmountOut.String(), tf.SandwichTrade.DiffTradeTokenBalance.String())
		actualDiff := new(big.Int).Sub(tf.SandwichTrade.AmountOut, tf.SandwichTrade.DiffTradeTokenBalance)
		log.Info().Msgf("actual diff %s", actualDiff.String())
		percentDiff := artemis_eth_units.PercentDiff(tf.SandwichTrade.AmountOut, tf.SandwichTrade.DiffTradeTokenBalance)
		log.Info().Msgf("percent diff %s", percentDiff.String())
		return nil, errors.New("amount out is less than the diff trade token balance")
	}
	if tf.SandwichTrade.AmountOut.String() == "0" {
		log.Info().Msgf("sandwich amount out is 0")
		return nil, errors.New("sandwich amount out is 0")
	}

	tf.SandwichTrade.PostTradeEthBalance = endEthBal
	tf.SandwichTrade.DiffTradeEthBalance = new(big.Int).Sub(endEthBal, startEthBal)
	return nil, nil
}

func (u *UniswapClient) SandwichTradeGetAmountsOut(tf *TradeExecutionFlow) ([]*big.Int, error) {
	pathSlice := []string{tf.SandwichTrade.AmountInAddr.String(), tf.SandwichTrade.AmountOutAddr.String()}
	amountsOutFirstPair, err := u.GetAmountsOut(tf.Tx.To(), tf.SandwichTrade.AmountIn, pathSlice)
	if err != nil {
		return nil, err
	}
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

	if artemis_eth_units.IsXLessThanY(tf.SandwichTrade.AmountOut, amountsOutFirstPair[1]) {
		log.Warn().Msgf(fmt.Sprintf("amount out not equal to expected amount out %s, actual amount out: %s", tf.UserTrade.AmountOut.String(), amountsOutFirstPair[1].String()))
		percentDiff := artemis_eth_units.PercentDiff(tf.SandwichTrade.AmountOut, amountsOutFirstPair[1])
		log.Info().Msgf("percent diff %s", percentDiff.String())
		return amountsOutFirstPair, errors.New("amount out not equal to expected")
	}
	tf.SandwichTrade.SimulatedAmountOut = amountsOutFirstPair[1]
	return amountsOutFirstPair, err
}
