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
	// NOTE: permissive transfer tax during historical sim
	tf.SandwichTrade.AmountOut = artemis_eth_units.ApplyTransferTax(tf.SandwichTrade.SimulatedAmountOut, 3, 1000)
	bal, err := u.Web3Client.ReadERC20TokenBalance(ctx, tf.SandwichTrade.AmountOutAddr.String(), u.Web3Client.PublicKey())
	if err != nil {
		log.Err(err).Msg("error getting pre trade amount out token balance")
		return nil, err
	}
	tf.SandwichTrade.PreTradeTokenBalance = bal
	fmt.Println("pre trade amount out token balance", bal.String())

	tf.SandwichTrade.PreTradeEthBalance = startEthBal
	if tf.InitialPairV3 != nil {
		err = u.ExecTradeV3SwapFromTokenToToken(ctx, tf.InitialPairV3, &tf.SandwichTrade)
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

	if !artemis_eth_units.IsXGreaterThanOrEqualToY(tf.SandwichTrade.DiffTradeTokenBalance, tf.SandwichTrade.AmountOut) {
		if !artemis_eth_units.PercentDiffFloatComparison(tf.SandwichTrade.AmountOut, tf.SandwichTrade.DiffTradeTokenBalance, 0.01) {
			log.Info().Msgf("sandwich trade: amount out %s is less than the diff trade token balance %s", tf.SandwichTrade.AmountOut.String(), tf.SandwichTrade.DiffTradeTokenBalance.String())
			actualDiff := new(big.Int).Sub(tf.SandwichTrade.AmountOut, tf.SandwichTrade.DiffTradeTokenBalance)
			log.Info().Msgf("sandwich trade: actual diff %s", actualDiff.String())
			percentDiff := artemis_eth_units.PercentDiffFloat(tf.SandwichTrade.AmountIn, tf.SandwichTrade.DiffTradeTokenBalance)
			log.Info().Msgf("sandwich trade: percent diff %f", percentDiff)
			return nil, errors.New("sandwich trade: amount out is less than the diff trade token balance")
		}
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

	//amountOut := amountsOutFirstPair[1]
	//amountOut, err = artemis_pricing_utils.ApplyTransferTax(accounts.HexToAddress(tf.SandwichTrade.AmountOutAddr.String()), amountOut)
	//if err != nil {
	//	return nil, err
	//}
	//amountOut, err = artemis_pricing_utils.ApplyTransferTax(accounts.HexToAddress(tf.SandwichTrade.AmountInAddr.String()), amountOut)
	//if err != nil {
	//	return nil, err
	//}
	//amountOut = artemis_eth_units.SetSlippage(amountOut)
	//amountsOutFirstPair[1] = amountOut

	if u.DebugPrint {
		fmt.Println("sandwich trade trade path", pathSlice[0], pathSlice[1])
		fmt.Println("sandwich trade expected amount in", tf.SandwichTrade.AmountIn.String(), "amount out", tf.SandwichTrade.AmountOut.String())
		fmt.Println("sandwich trade simulated amount in", amountsOutFirstPair[0].String(), "amount out", amountsOutFirstPair[1].String())
	}

	if !artemis_eth_units.PercentDiffFloatComparison(tf.SandwichTrade.AmountIn, amountsOutFirstPair[0], 0.3) {
		artemis_eth_units.PercentDiffHighPrecision(tf.SandwichTrade.AmountIn, amountsOutFirstPair[0])
		log.Warn().Msgf(fmt.Sprintf("sandwich trade: amount in not equal to expected amount in %s, actual amount in: %s", tf.SandwichTrade.AmountIn.String(), amountsOutFirstPair[0].String()))
		return amountsOutFirstPair, errors.New("sandwich trade: amount in not equal to expected")
	}

	// 	if !artemis_eth_units.PercentDiffFloatComparison(tf.SandwichTrade.AmountOut, amountsOutFirstPair[1], 0.0001) {
	if !artemis_eth_units.IsXGreaterThanOrEqualToY(amountsOutFirstPair[1], tf.SandwichTrade.AmountOut) {
		if !artemis_eth_units.PercentDiffFloatComparison(tf.SandwichTrade.AmountOut, amountsOutFirstPair[1], 0.3) {
			log.Warn().Msgf(fmt.Sprintf("sandwich trade: amount out not equal to expected amount out %s, actual amount out: %s", tf.SandwichTrade.AmountOut.String(), amountsOutFirstPair[1].String()))
			percentDiff := artemis_eth_units.PercentDiffFloat(tf.SandwichTrade.AmountOut, amountsOutFirstPair[1])
			log.Info().Msgf("sandwich trade: percent diff %f", percentDiff)
			return amountsOutFirstPair, errors.New("sandwich trade: amount out not equal to expected")
		}
	}

	tf.SandwichTrade.SimulatedAmountOut = amountsOutFirstPair[1]
	return amountsOutFirstPair, err
}
