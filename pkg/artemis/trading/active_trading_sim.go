package artemis_realtime_trading

import (
	"context"
	"errors"

	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) SimToPackageTxBundle(ctx context.Context, tf *web3_client.TradeExecutionFlow, bypassSim bool) error {
	bundle := &artemis_flashbots.MevTxBundle{}
	if !bypassSim {
		// TODO set hardhat to live network
		err := a.a.U.Web3Client.MatchFrontRunTradeValues(tf)
		if err != nil {
			return err
		}
		err = a.InitActiveTradingSimEnv(ctx, tf)
		if err != nil {
			return err
		}
	}
	// FRONT_RUN
	if tf.InitialPairV3 != nil {
		//err := a.a.U.ExecTradeV3SwapFromTokenToToken(ctx, tf.InitialPairV3, &tf.FrontRunTrade)
		//if err != nil {
		//	return err
		//}
		return errors.New("uniswap V3 not supported yet")
	} else {
		//err := a.ExecTradeV2SwapFromTokenToToken(ctx, &tf.FrontRunTrade, bypassSim)
		//if err != nil {
		//	return err
		//}
	}
	bundle.AddTxs(tf.Tx)
	// FRONT_RUN

	// USER TRADE
	if !bypassSim {
		err := a.a.U.Web3Client.SendSignedTransaction(ctx, tf.Tx)
		if err != nil {
			return err
		}
	}
	bundle.AddTxs(tf.Tx)
	// USER TRADE

	// SANDWICH TRADE
	if tf.InitialPairV3 != nil {
		//err = a.a.U.ExecTradeV3SwapFromTokenToToken(ctx, tf.InitialPairV3, &tf.SandwichTrade)
		//if err != nil {
		//	return err
		//}
		return errors.New("uniswap V3 not supported yet")
	} else {
		//err := a.ExecTradeV2SwapFromTokenToToken(ctx, &tf.SandwichTrade, bypassSim)
		//if err != nil {
		//	return err
		//}
	}
	bundle.AddTxs(tf.Tx)
	tf.Bundle = bundle
	return nil
}

//err = tf.GetAggregateGasUsage(ctx, U.Web3Client)
//if err != nil {
//	U.TradeFailureReport.EndStage = "post trade getting gas usage"
//	log.Err(err).Msg("error getting aggregate gas usage")
//	return U.MarkEndOfSimDueToErr(err)
//}
//err = U.VerifyTradeResults(tf)
//if err != nil {
//	U.TradeFailureReport.EndStage = "verifying trade results"
//	log.Err(err).Msg("error verifying trade results")
//	return U.MarkEndOfSimDueToErr(err)
//}
