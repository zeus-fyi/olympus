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
		fmt.Println("user trade expected amount in", tf.UserTrade.AmountIn.String(), "amount out", tf.UserTrade.AmountOut.String())
		fmt.Println("user trade simulated amount in", amountsOutFirstPair[0].String(), "amount out", amountsOutFirstPair[1].String())
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
