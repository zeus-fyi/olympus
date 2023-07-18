package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) PackageSandwich(ctx context.Context, tf *web3_client.TradeExecutionFlow) (*flashbotsrpc.FlashbotsSendBundleResponse, error) {
	if tf == nil || tf.Tx == nil {
		return nil, errors.New("tf is nil")
	}
	// front run
	ur, err := a.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.FrontRunTrade)
	if err != nil {
		return nil, err
	}
	tx, err := a.universalRouterCmdToTxBuilder(ctx, ur)
	if err != nil {
		log.Err(err).Interface("txHash", tx.Hash().String()).Msg("failed to add tx to bundle group")
		return nil, err
	}
	// user trade
	mevTx, err := a.AddTxToBundleGroup(ctx, tf.Tx)
	if err != nil {
		log.Err(err).Interface("mevTx", mevTx).Msg("failed to add tx to bundle group")
		return nil, err
	}
	// sandwich trade
	ur, err = a.GenerateTradeV2SwapFromTokenToToken(ctx, ur, &tf.SandwichTrade)
	if err != nil {
		log.Err(err).Msg("failed to generate sandwich tx")
		return nil, err
	}
	txSand, err := a.universalRouterCmdToTxBuilder(ctx, ur)
	if err != nil {
		log.Err(err).Interface("txSand", txSand.Hash().String()).Msg("failed to add tx to bundle group")
		return nil, err
	}
	return nil, err
}

func (a *AuxiliaryTradingUtils) StagingPackageSandwichAndCall(ctx context.Context, tf *web3_client.TradeExecutionFlow) (*flashbotsrpc.FlashbotsCallBundleResponse, error) {
	_, err := a.PackageSandwich(ctx, tf)
	if err != nil {
		log.Err(err).Msg("failed to package sandwich")
		return nil, err
	}
	resp, err := a.CallFlashbotsBundleStaging(ctx, a.MevTxGroup.MevTxs)
	if err != nil {
		log.Err(err).Interface("fbCallResp", resp).Msg("failed to send sandwich")
		return nil, err
	}
	log.Info().Interface("fbCallResp", resp).Msg("sent sandwich")
	return &resp, err
}

func (a *AuxiliaryTradingUtils) PackageSandwichAndSend(ctx context.Context, tf *web3_client.TradeExecutionFlow) (*flashbotsrpc.FlashbotsSendBundleResponse, error) {
	_, err := a.PackageSandwich(ctx, tf)
	if err != nil {
		log.Err(err).Msg("failed to package sandwich")
		return nil, err
	}
	_, err = a.CallAndSendFlashbotsBundle(ctx)
	if err != nil {
		log.Err(err).Msg("failed to send sandwich")
		return nil, err
	}
	return nil, err
}
