package artemis_realtime_trading

import (
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) SimW3c() *web3_client.Web3Client {
	return &a.us.a.U.Web3Client
}
func (a *ActiveTrading) GetSimUniswapClient() *web3_client.UniswapClient {
	return a.us.a.U
}
func (a *ActiveTrading) GetSimAuxClient() *artemis_trading_auxiliary.AuxiliaryTradingUtils {
	return a.us.a
}

//func (a *ActiveTrading) SimStage(ctx context.Context, w3c web3_client.Web3Client, tfSlice []web3_client.TradeExecutionFlowJSON) ErrWrapper {
//	err := a.SimTxFilter(ctx, tfSlice)
//	if err != nil {
//		log.Err(err).Msg("failed to pass sim tx filter")
//		return ErrWrapper{Err: err, Stage: "SimTxFilter"}
//	}
//	if len(tfSlice) <= 0 {
//		return ErrWrapper{Err: errors.New("SimTxFilter: no tx flows to simulate"), Stage: "SimTxFilter"}
//	}
//	a.GetMetricsClient().StageProgressionMetrics.CountPostSimFilterTx(float64(1))
//	log.Info().Msg("starting simulation")
//	err = a.SimToPackageTxBundles(ctx, tfSlice, false)
//	if err != nil {
//		log.Err(err).Msg("failed to simulate txs")
//		return ErrWrapper{Err: err, Stage: "SimToPackageTxBundles", Code: 200}
//	}
//	log.Info().Msg("simulation stage complete: starting active trading filter")
//	a.GetMetricsClient().StageProgressionMetrics.CountPostSimStage(float64(len(tfSlice)))
//	err = ActiveTradingFilterSlice(ctx, w3c, tfSlice)
//	if err != nil {
//		log.Err(err).Msg("failed to pass active trading filter")
//		return ErrWrapper{Err: err, Stage: "ActiveTradingFilterSlice", Code: 200}
//	}
//	log.Info().Msg("preparing bundles for submission")
//	a.GetMetricsClient().StageProgressionMetrics.CountPostActiveTradingFilter(float64(len(tfSlice)))
//	err = artemis_realtime_trading.ProcessBundleStage(ctx, w3c, tfSlice, a.GetMetricsClient())
//	if err != nil {
//		log.Err(err).Msg("failed to process bundles")
//		return ErrWrapper{Err: err, Stage: "ProcessBundleStage", Code: 200}
//	}
//	log.Info().Msg("bundles successfully sent")
//	a.GetMetricsClient().StageProgressionMetrics.CountSentFlashbotsBundleSubmission(float64(len(tfSlice)))
//	return ErrWrapper{Err: nil, Stage: "Success", Code: 200}
//}

//func (a *ActiveTrading) SimToPackageTxBundles(ctx context.Context, tfSlide []web3_client.TradeExecutionFlowJSON, bypassSim bool) error {
//	for _, tf := range tfSlide {
//		tfConv := tf.ConvertToBigIntType()
//		if tf.FrontRunTrade.AmountInAddr.String() != artemis_trading_constants.WETH9ContractAddressAccount.String() {
//			return errors.New("SimToPackageTxBundles: profit token is not WETH")
//		}
//		err := a.SimToPackageTxBundle(ctx, &tfConv, bypassSim)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

//func (a *ActiveTrading) SimToPackageTxBundle(ctx context.Context, tf *web3_client.TradeExecutionFlow, bypassSim bool) error {
//	log.Info().Msg("SimToPackageTxBundle: starting")
//	if tf == nil {
//		return errors.New("tf is nil")
//	}
//	bundle := &artemis_flashbots.MevTxBundle{}
//	if !bypassSim {
//		a.GetSimAuxClient().U.Web3Client.AddSessionLockHeader(tf.Tx.Hash().String())
//		// TODO set hardhat to live network
//		err := a.setupCleanSimEnvironment(ctx, tf)
//		if err != nil {
//			log.Err(err).Msg("SimToPackageTxBundle: failed to setup clean sim environment")
//			return err
//		}
//		err = a.SimW3c().MatchFrontRunTradeValues(tf)
//		if err != nil {
//			log.Err(err).Msg("SimToPackageTxBundle: failed to match front run trade values")
//			return err
//		}
//	}
//	// FRONT_RUN
//	if tf.InitialPairV3 != nil {
//		//err := a.a.U.ExecTradeV3SwapFromTokenToToken(ctx, tf.InitialPairV3, &tf.FrontRunTrade)
//		//if err != nil {
//		//	return err
//		//}
//		return errors.New("uniswap V3 not supported yet")
//	} else {
//
//		start := tf.FrontRunTrade.AmountOut
//		num := 0
//		denom := 1000
//		for i := 1; i < 7; i++ {
//			switch i {
//			case 0:
//				num = 1
//				denom = 1
//			case 1:
//				num = 1
//				denom = 1000
//			case 2:
//				num = 10
//				denom = 1000
//			case 3:
//				num = 50
//				denom = 1000
//			case 4:
//				num = 100
//				denom = 1000
//			case 5:
//				num = 200
//				denom = 1000
//			default:
//				return errors.New("failed to find a valid transfer tax")
//			}
//			tf.FrontRunTrade.AmountOut = artemis_eth_units.ApplyTransferTax(start, num, denom)
//			if tf.FrontRunTrade.AmountOut.String() == "0" {
//				return errors.New("amount out was set to zero")
//			}
//			ur, _, err := a.GetSimAuxClient().GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.FrontRunTrade)
//			if err != nil {
//				fmt.Println("failed to generate trade", ur.Commands)
//				log.Err(err).Msg("SimToPackageTxBundle: failed to generate trade")
//				return err
//			}
//			_, err = a.GetSimAuxClient().U.ExecUniswapUniversalRouterCmd(*ur)
//			if err == nil {
//				log.Info().Interface("num", num).Interface("denom", denom).Msgf("Injected trade with amount out: %s", tf.FrontRunTrade.AmountOut.String())
//				break
//			}
//		}
//	}
//	err := bundle.AddTxs(tf.FrontRunTrade.BundleTxs...)
//	if err != nil {
//		return err
//	}
//	// FRONT_RUN
//
//	// USER TRADE
//	if !bypassSim {
//		err = a.SimW3c().SendImpersonatedTx(ctx, tf.Tx)
//		if err != nil {
//			log.Err(err).Msg("SimToPackageTxBundle: failed to send impersonated tx")
//			return err
//		}
//	}
//	err = bundle.AddTxs(tf.Tx)
//	if err != nil {
//		return err
//	}
//	// USER TRADE
//
//	// SANDWICH TRADE
//	if tf.InitialPairV3 != nil {
//		//err = a.a.U.ExecTradeV3SwapFromTokenToToken(ctx, tf.InitialPairV3, &tf.SandwichTrade)
//		//if err != nil {
//		//	return err
//		//}
//		return errors.New("uniswap V3 not supported yet")
//	} else {
//		tf.SandwichTrade.AmountIn = tf.FrontRunTrade.AmountOut
//		_, err = a.GetSimAuxClient().U.SandwichTradeGetAmountsOut(tf)
//		if err != nil {
//			return err
//		}
//		tf.SandwichTrade.AmountOut = tf.SandwichTrade.SimulatedAmountOut
//		ur, _, serr := a.GetSimAuxClient().GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.SandwichTrade)
//		if serr != nil {
//			fmt.Println("SimToPackageTxBundle: failed to generate trade", ur.Commands)
//			log.Err(serr).Msg("failed to generate trade")
//			return serr
//		}
//	}
//	err = bundle.AddTxs(tf.SandwichTrade.BundleTxs...)
//	if err != nil {
//		return err
//	}
//	tf.Bundle = bundle
//	return nil
//}
