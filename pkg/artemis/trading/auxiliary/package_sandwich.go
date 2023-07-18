package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) PackageSandwich(ctx context.Context, tf *web3_client.TradeExecutionFlow) (*MevTxGroup, error) {
	if tf == nil || tf.Tx == nil {
		return nil, errors.New("tf is nil")
	}
	bundle := &MevTxGroup{
		EventID:    0,
		OrderedTxs: []TxWithMetadata{},
		MevTxs:     []artemis_eth_txs.EthTx{},
	}
	startCtx := ctx
	// front run
	frontRunCtx := a.CreateFrontRunCtx(startCtx)
	ur, fpt, err := a.GenerateTradeV2SwapFromTokenToToken(frontRunCtx, nil, &tf.FrontRunTrade)
	if err != nil {
		log.Err(err).Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: failed to generate front run tx")
		return nil, err
	}
	frontRunTx, err := a.universalRouterCmdToTxBuilder(frontRunCtx, ur)
	if err != nil {
		log.Err(err).Interface("txHash", frontRunTx.Hash().String()).Msg("FRONT_RUN: failed to add tx to bundle group")
		return nil, err
	}
	frTx := TxWithMetadata{
		Tx: frontRunTx,
	}
	if fpt != nil {
		frTx.Permit2Tx = fpt.Permit2Tx
	}
	bundle, err = a.AddTxToBundleGroup(frontRunCtx, frTx, bundle)
	if err != nil {
		log.Info().Interface("mevTx", bundle.MevTxs).Msg("FRONT_RUN: error adding tx to bundle group")
		return nil, err
	}
	// user trade
	userCtx := a.CreateUserTradeCtx(startCtx)
	userTx := TxWithMetadata{
		Tx: tf.Tx,
	}
	bundle, err = a.AddTxToBundleGroup(userCtx, userTx, bundle)
	if err != nil {
		log.Err(err).Interface("mevTx", bundle.MevTxs).Msg("USER_TRADE: failed to add tx to bundle group")
		return nil, err
	}
	// sandwich trade
	backRunCtx := a.CreateBackRunCtx(startCtx)
	ur, spt, err := a.GenerateTradeV2SwapFromTokenToToken(backRunCtx, ur, &tf.SandwichTrade)
	if err != nil {
		log.Err(err).Msg("SANDWICH_TRADE: failed to generate sandwich tx")
		return nil, err
	}
	txSand, err := a.universalRouterCmdToTxBuilder(backRunCtx, ur)
	if err != nil {
		log.Err(err).Interface("txSand", txSand.Hash().String()).Msg("SANDWICH_TRADE: failed to add tx to bundle group")
		return nil, err
	}
	sandwichTx := TxWithMetadata{
		Tx: txSand,
	}
	if spt != nil {
		sandwichTx.Permit2Tx = spt.Permit2Tx
	}
	bundle, err = a.AddTxToBundleGroup(backRunCtx, sandwichTx, bundle)
	if err != nil {
		log.Err(err).Interface("mevTx", bundle.MevTxs).Msg("SANDWICH_TRADE: failed to add tx to bundle group")
		return nil, err
	}
	if len(bundle.MevTxs) != 3 {
		log.Warn().Int("bundleTxCount", len(bundle.MevTxs)).Msg("SANDWICH_TRADE: sandwich bundle not 3 txs")
		return nil, errors.New("sandwich bundle not 3 txs")
	}
	return bundle, err
}

func (a *AuxiliaryTradingUtils) StagingPackageSandwichAndCall(ctx context.Context, tf *web3_client.TradeExecutionFlow) (*flashbotsrpc.FlashbotsCallBundleResponse, error) {
	bundle, err := a.PackageSandwich(ctx, tf)
	if err != nil {
		log.Err(err).Msg("failed to package sandwich")
		return nil, err
	}
	if bundle == nil {
		return nil, errors.New("bundle is nil")
	}
	resp, err := a.CallFlashbotsBundleStaging(ctx, *bundle)
	if err != nil {
		log.Err(err).Interface("fbCallResp", resp).Msg("failed to send sandwich")
		return nil, err
	}
	log.Info().Interface("fbCallResp", resp).Msg("sent sandwich")
	return &resp, err
}

func (a *AuxiliaryTradingUtils) PackageSandwichAndSend(ctx context.Context, tf *web3_client.TradeExecutionFlow) (*flashbotsrpc.FlashbotsSendBundleResponse, error) {
	bundle, err := a.PackageSandwich(ctx, tf)
	if err != nil {
		log.Err(err).Msg("failed to package sandwich")
		return nil, err
	}
	if bundle == nil {
		return nil, errors.New("bundle is nil")
	}
	_, err = a.CallAndSendFlashbotsBundle(ctx, *bundle)
	if err != nil {
		log.Err(err).Msg("failed to send sandwich")
		return nil, err
	}
	return nil, err
}
