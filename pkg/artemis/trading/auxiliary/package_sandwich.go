package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func PackageSandwich(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow) (*MevTxGroup, error) {
	log.Info().Msg("PackageSandwich: start")
	if tf == nil || tf.Tx == nil {
		return nil, errors.New("tf is nil")
	}
	bundle := &MevTxGroup{
		EventID:      0,
		OrderedTxs:   []TxWithMetadata{},
		MevTxs:       []artemis_eth_txs.EthTx{},
		TotalGasCost: artemis_eth_units.NewBigInt(0),
	}
	startCtx := ctx
	// front run
	frontRunCtx := CreateFrontRunCtx(startCtx)
	ur, fpt, err := GenerateTradeV2SwapFromTokenToToken(frontRunCtx, w3c, nil, &tf.FrontRunTrade)
	if err != nil {
		log.Err(err).Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: failed to generate front run tx")
		return nil, err
	}
	frontRunTx, scInfoFrontRun, err := universalRouterCmdToTxBuilder(frontRunCtx, w3c, ur)
	if err != nil {
		log.Err(err).Interface("txHash", frontRunTx.Hash().String()).Msg("FRONT_RUN: failed to add tx to bundle group")
		return nil, err
	}

	frontRunGasCost := artemis_eth_units.MulBigInt(artemis_eth_units.AddBigInt(scInfoFrontRun.GasFeeCap, scInfoFrontRun.GasTipCap), artemis_eth_units.NewBigIntFromUint(scInfoFrontRun.GasLimit))
	tf.FrontRunTrade.TotalGasCost = frontRunGasCost.Uint64()
	bundle.TotalGasCost = artemis_eth_units.AddBigInt(bundle.TotalGasCost, frontRunGasCost)

	frTx := TxWithMetadata{
		Tx: frontRunTx,
	}
	if fpt != nil {
		frTx.Permit2Tx = fpt.Permit2Tx
	}
	bundle, err = AddTxToBundleGroup(frontRunCtx, frTx, bundle)
	if err != nil {
		log.Info().Interface("mevTx", bundle.MevTxs).Msg("FRONT_RUN: error adding tx to bundle group")
		return nil, err
	}
	// user trade
	userCtx := CreateUserTradeCtx(startCtx)
	userTx := TxWithMetadata{
		Tx: tf.Tx,
	}
	bundle, err = AddTxToBundleGroup(userCtx, userTx, bundle)
	if err != nil {
		log.Err(err).Interface("mevTx", bundle.MevTxs).Msg("USER_TRADE: failed to add tx to bundle group")
		return nil, err
	}
	// sandwich trade
	backRunCtx := CreateBackRunCtx(startCtx, w3c)
	ur, spt, err := GenerateTradeV2SwapFromTokenToToken(backRunCtx, w3c, ur, &tf.SandwichTrade)
	if err != nil {
		log.Err(err).Msg("SANDWICH_TRADE: failed to generate sandwich tx")
		return nil, err
	}

	txSand, scInfoSand, err := universalRouterCmdToTxBuilder(backRunCtx, w3c, ur)
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
	sandwichGasCost := artemis_eth_units.MulBigInt(artemis_eth_units.AddBigInt(scInfoSand.GasFeeCap, scInfoSand.GasTipCap), artemis_eth_units.NewBigIntFromUint(scInfoSand.GasLimit))
	tf.SandwichTrade.TotalGasCost = sandwichGasCost.Uint64()
	bundle.TotalGasCost = artemis_eth_units.AddBigInt(bundle.TotalGasCost, sandwichGasCost)
	bundle, err = AddTxToBundleGroup(backRunCtx, sandwichTx, bundle)
	if err != nil {
		log.Err(err).Interface("mevTx", bundle.MevTxs).Msg("SANDWICH_TRADE: failed to add tx to bundle group")
		return nil, err
	}
	if len(bundle.MevTxs) != 3 {
		log.Warn().Int("bundleTxCount", len(bundle.MevTxs)).Msg("SANDWICH_TRADE: sandwich bundle not 3 txs")
		return nil, errors.New("sandwich bundle not 3 txs")
	}
	log.Info().Msg("PackageSandwich: end")
	return bundle, err
}

func (a *AuxiliaryTradingUtils) StagingPackageSandwichAndCall(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow) (*flashbotsrpc.FlashbotsCallBundleResponse, *MevTxGroup, error) {
	log.Info().Msg("StagingPackageSandwichAndCall: start")
	if tf == nil {
		return nil, nil, errors.New("tf is nil")
	}
	bundle, err := PackageSandwich(ctx, w3c, tf)
	if err != nil {
		log.Err(err).Msg("StagingPackageSandwichAndCall: failed to package sandwich")
		return nil, nil, err
	}
	if bundle == nil {
		return nil, nil, errors.New("bundle is nil")
	}
	//log.Info().Interface("bundle", bundle).Msg("isBundleProfitHigherThanGasFee: bundle")
	//ok, err := isBundleProfitHigherThanGasFee(bundle, tf)
	//if err != nil {
	//	log.Err(err).Bool("ok", ok).Msg("StagingPackageSandwichAndCall: isBundleProfitHigherThanGasFee: failed to check if profit is higher than gas fee")
	//	return nil, nil, err
	//}
	resp, err := a.CallFlashbotsBundleStaging(ctx, w3c, *bundle)
	if err != nil {
		log.Err(err).Interface("fbCallResp", resp).Msg("failed to send sandwich")
		return nil, nil, err
	}
	log.Info().Interface("fbCallResp", resp).Msg("sent sandwich")
	return &resp, bundle, err
}

func (a *AuxiliaryTradingUtils) PackageSandwichAndSend(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow) (*flashbotsrpc.FlashbotsSendBundleResponse, error) {
	bundle, err := PackageSandwich(ctx, w3c, tf)
	if err != nil {
		log.Err(err).Msg("failed to package sandwich")
		return nil, err
	}
	if bundle == nil {
		return nil, errors.New("bundle is nil")
	}
	_, err = a.CallAndSendFlashbotsBundle(ctx, w3c, *bundle)
	if err != nil {
		log.Err(err).Msg("failed to send sandwich")
		return nil, err
	}
	return nil, err
}
