package artemis_realtime_trading

import (
	"context"

	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/flashbots"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) SimulateTx(ctx context.Context, tf *web3_client.TradeExecutionFlow) error {
	// TODO set hardhat to live network
	var bundle *artemis_flashbots.MevTxBundle
	err := a.u.Web3Client.MatchFrontRunTradeValues(tf)
	if err != nil {
		return err
	}

	err = a.InitActiveTradingSimEnv(ctx, tf)
	if err != nil {
		return err
	}
	// FRONT_RUN
	if tf.InitialPairV3 != nil {
		err = a.u.ExecTradeV3SwapFromTokenToToken(ctx, tf.InitialPairV3, &tf.FrontRunTrade)
		if err != nil {
			return err
		}
	} else {
		err = a.u.ExecTradeV2SwapFromTokenToToken(ctx, &tf.FrontRunTrade)
		if err != nil {
			return err
		}
	}
	bundle.Txs = append(bundle.Txs, tf.FrontRunTrade.BundleTxs...)
	// USER TRADE
	err = a.u.Web3Client.SendSignedTransaction(ctx, tf.Tx)
	if err != nil {
		return err
	}
	bundle.Txs = append(bundle.Txs, tf.Tx)
	// SANDWICH TRADE
	if tf.InitialPairV3 != nil {
		err = a.u.ExecTradeV3SwapFromTokenToToken(ctx, tf.InitialPairV3, &tf.SandwichTrade)
		if err != nil {
			return err
		}
	} else {
		err = a.u.ExecTradeV2SwapFromTokenToToken(ctx, &tf.SandwichTrade)
		if err != nil {
			return err
		}
	}
	bundle.Txs = append(bundle.Txs, tf.SandwichTrade.BundleTxs...)
	//err = tf.GetAggregateGasUsage(ctx, u.Web3Client)
	//if err != nil {
	//	u.TradeFailureReport.EndStage = "post trade getting gas usage"
	//	log.Err(err).Msg("error getting aggregate gas usage")
	//	return u.MarkEndOfSimDueToErr(err)
	//}
	//err = u.VerifyTradeResults(tf)
	//if err != nil {
	//	u.TradeFailureReport.EndStage = "verifying trade results"
	//	log.Err(err).Msg("error verifying trade results")
	//	return u.MarkEndOfSimDueToErr(err)
	//}
	//if !u.TestMode {
	//	return u.MarkEndOfSimDueToErr(nil)
	//}
	tf.Bundle = bundle
	return nil
}
