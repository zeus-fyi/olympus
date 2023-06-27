package web3_client

import (
	"context"
	"math/big"

	"github.com/rs/zerolog/log"
)

/*
	part 0. trade filter

  - should filter trades on multiple criteria
	- eg. profit denomination etc
	- eg. token risk score

	part 1. fast parallel processing of txs

  - should decode tx
  - should set up balances
  - should minimize rpc calls

	part 2. bundle txs into a single tx

  - should bundle txs into a single tx from processed ones

Overall:
	- should capture sim runtime metrics
	- should send trade metrics to prometheus
*/

func (u *UniswapClient) ActiveTradingRealTimeTxAnalysis(ctx context.Context, tf *TradeExecutionFlow) error {
	u.TradeAnalysisReport.TradeMethod = tf.Trade.TradeMethod
	u.TradeAnalysisReport.AmountInAddr = tf.FrontRunTrade.AmountInAddr.String()
	u.TradeAnalysisReport.AmountOutAddr = tf.SandwichTrade.AmountOutAddr.String()
	// this isn't included in trade gas costs since we amortize one time gas costs for permit2
	max, _ := new(big.Int).SetString(maxUINT, 10)
	approveTx, err := u.ApproveSpender(ctx, WETH9ContractAddress, Permit2SmartContractAddress, max)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}

	secondToken := tf.FrontRunTrade.AmountInAddr.String()
	if tf.FrontRunTrade.AmountInAddr.String() == WETH9ContractAddress {
		secondToken = tf.FrontRunTrade.AmountOutAddr.String()
	}
	approveTx, err = u.ApproveSpender(ctx, secondToken, Permit2SmartContractAddress, max)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}

	err = u.Web3Client.MatchFrontRunTradeValues(tf)
	if err != nil {
		u.TradeFailureReport.EndStage = "executing front run balance setup"
		log.Err(err).Msg("error executing front run balance setup")
		return u.MarkEndOfSimDueToErr(err)
	}
	_, err = u.ExecFrontRunTradeStepTokenTransfer(tf)
	if err != nil {
		u.TradeFailureReport.EndStage = "executing front run trade"
		log.Err(err).Msg("error executing front run trade step token transfer")
		return u.MarkEndOfSimDueToErr(err)
	}
	_, err = u.ExecUserTradeStep(tf)
	if err != nil {
		u.TradeFailureReport.EndStage = "executing user trade step"
		log.Err(err).Msg("error executing user trade step")
		return u.MarkEndOfSimDueToErr(err)
	}
	_, err = u.ExecSandwichTradeStepTokenTransfer(tf)
	if err != nil {
		u.TradeFailureReport.EndStage = "executing sandwich trade"
		log.Err(err).Msg("error executing sandwich trade step token transfer")
		return u.MarkEndOfSimDueToErr(err)
	}
	err = tf.GetAggregateGasUsage(ctx, u.Web3Client)
	if err != nil {
		u.TradeFailureReport.EndStage = "post trade getting gas usage"
		log.Err(err).Msg("error getting aggregate gas usage")
		return u.MarkEndOfSimDueToErr(err)
	}
	err = u.VerifyTradeResults(tf)
	if err != nil {
		u.TradeFailureReport.EndStage = "verifying trade results"
		log.Err(err).Msg("error verifying trade results")
		return u.MarkEndOfSimDueToErr(err)
	}
	if !u.TestMode {
		return u.MarkEndOfSimDueToErr(nil)
	}
	return nil
}
