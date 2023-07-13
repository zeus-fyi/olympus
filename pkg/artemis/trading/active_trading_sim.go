package artemis_realtime_trading

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) simW3c() *web3_client.Web3Client {
	return &a.us.a.U.Web3Client
}
func (a *ActiveTrading) simAuxUtils() *artemis_trading_auxiliary.AuxiliaryTradingUtils {
	return a.us.a
}

func (a *ActiveTrading) SimToPackageTxBundles(ctx context.Context, tfSlide []web3_client.TradeExecutionFlowJSON, bypassSim bool) error {
	for _, tf := range tfSlide {
		tfConv := tf.ConvertToBigIntType()
		if tf.FrontRunTrade.AmountInAddr.String() != artemis_trading_constants.WETH9ContractAddressAccount.String() {
			return errors.New("profit token is not WETH")
		}
		err := a.SimToPackageTxBundle(ctx, &tfConv, bypassSim)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ActiveTrading) SimToPackageTxBundle(ctx context.Context, tf *web3_client.TradeExecutionFlow, bypassSim bool) error {
	if tf == nil {
		return errors.New("tf is nil")
	}
	bundle := &artemis_flashbots.MevTxBundle{}
	if !bypassSim {
		a.simAuxUtils().U.Web3Client.AddSessionLockHeader(tf.Tx.Hash().String())
		// TODO set hardhat to live network
		err := a.setupCleanSimEnvironment(ctx, tf)
		if err != nil {
			log.Err(err).Msg("failed to setup clean sim environment")
			return err
		}
		err = a.simW3c().MatchFrontRunTradeValues(tf)
		if err != nil {
			log.Err(err).Msg("failed to match front run trade values")
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
		ur, err := a.simAuxUtils().GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.FrontRunTrade)
		if err != nil {
			fmt.Println("failed to generate trade", ur.Commands)
			log.Err(err).Msg("failed to generate trade")
			return err
		}
		_, err = a.simAuxUtils().U.ExecUniswapUniversalRouterCmd(*ur)
		if err != nil {
			log.Err(err).Msg("failed to execute trade")
			return err
		}

	}
	err := bundle.AddTxs(tf.FrontRunTrade.BundleTxs...)
	if err != nil {
		return err
	}
	// FRONT_RUN

	// USER TRADE
	if !bypassSim {
		err = a.simW3c().SendImpersonatedTx(ctx, tf.Tx)
		if err != nil {
			log.Err(err).Msg("failed to send impersonated tx")
			return err
		}
	}
	err = bundle.AddTxs(tf.Tx)
	if err != nil {
		return err
	}
	// USER TRADE

	// SANDWICH TRADE
	if tf.InitialPairV3 != nil {
		//err = a.a.U.ExecTradeV3SwapFromTokenToToken(ctx, tf.InitialPairV3, &tf.SandwichTrade)
		//if err != nil {
		//	return err
		//}
		return errors.New("uniswap V3 not supported yet")
	} else {
		ur, serr := a.simAuxUtils().GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.SandwichTrade)
		if serr != nil {
			fmt.Println("failed to generate trade", ur.Commands)
			log.Err(serr).Msg("failed to generate trade")
			return serr
		}
	}
	err = bundle.AddTxs(tf.SandwichTrade.BundleTxs...)
	if err != nil {
		return err
	}
	tf.Bundle = bundle
	return nil
}
