package web3_client

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

func (u *UniswapClient) ExecUserTradeStep(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
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

func (u *UniswapClient) UserTradeGetAmountsOut(tf *TradeExecutionFlow) ([]*big.Int, error) {
	pathSlice := []string{tf.UserTrade.AmountInAddr.String(), tf.UserTrade.AmountOutAddr.String()}
	amountsOutFirstPair, err := u.GetAmountsOut(tf.Tx.To(), tf.UserTrade.AmountIn, pathSlice)
	if err != nil {
		return nil, err
	}
	if len(amountsOutFirstPair) != 2 {
		return nil, errors.New("amounts out not equal to expected")
	}
	if u.DebugPrint {
		fmt.Println("user trade trade path", pathSlice[0], pathSlice[1])
		fmt.Println("user trade expected amount in", tf.UserTrade.AmountIn.String(), "amount out", tf.UserTrade.AmountOut.String())
		fmt.Println("user trade simulated amount in", amountsOutFirstPair[0].String(), "amount out", amountsOutFirstPair[1].String())
	}

	if !artemis_eth_units.PercentDiffFloatComparison(tf.UserTrade.AmountIn, amountsOutFirstPair[0], 0.01) {
		log.Warn().Msgf(fmt.Sprintf("txhash %s: amount in not equal to expected amount in %s, actual amount in: %s", tf.Tx.Hash().String(), tf.UserTrade.AmountIn.String(), amountsOutFirstPair[0].String()))
		return amountsOutFirstPair, errors.New("amount in not equal to expected")
	}
	//if tf.UserTrade.AmountIn.String() != amountsOutFirstPair[0].String() {
	//	log.Warn().Msgf(fmt.Sprintf("txhash %s: amount in not equal to expected amount in %s, actual amount in: %s", tf.Tx.Hash().String(), tf.UserTrade.AmountIn.String(), amountsOutFirstPair[0].String()))
	//	return amountsOutFirstPair, errors.New("amount in not equal to expected")
	//}
	if !artemis_eth_units.IsXGreaterThanOrEqualToY(amountsOutFirstPair[1], tf.UserTrade.AmountOut) {
		if !artemis_eth_units.PercentDiffFloatComparison(tf.UserTrade.AmountOut, amountsOutFirstPair[1], 0.01) {
			log.Info().Msgf("user trade: amount out %s is less than the diff trade token balance %s", tf.UserTrade.AmountOut.String(), tf.UserTrade.DiffTradeTokenBalance.String())
			actualDiff := new(big.Int).Sub(tf.UserTrade.AmountOut, tf.UserTrade.DiffTradeTokenBalance)
			log.Info().Msgf("actual diff %s", actualDiff.String())
			percentDiff := artemis_eth_units.PercentDiffFloat(tf.UserTrade.AmountIn, tf.UserTrade.DiffTradeTokenBalance)
			log.Info().Msgf("percent diff %f", percentDiff)
		}
	}

	//if tf.UserTrade.AmountOut.String() != amountsOutFirstPair[1].String() {
	//	log.Warn().Msgf(fmt.Sprintf("txhash %s: amount out not equal to expected amount out %s, actual amount out: %s", tf.Tx.Hash().String(), tf.UserTrade.AmountOut.String(), amountsOutFirstPair[1].String()))
	//	if u.DebugPrint {
	//		diff := new(big.Int).Sub(amountsOutFirstPair[1], tf.UserTrade.AmountOut)
	//		fmt.Println("user trade actual - expected ", diff.String())
	//	}
	//	return amountsOutFirstPair, errors.New("amount out not equal to expected")
	//}
	tf.UserTrade.SimulatedAmountOut = amountsOutFirstPair[1]
	return amountsOutFirstPair, err
}
