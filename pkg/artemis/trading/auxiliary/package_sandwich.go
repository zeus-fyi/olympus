package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func packageFrontRun(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow) (*TxWithMetadata, error) {
	frontRunCtx := CreateFrontRunCtx(context.Background())
	ur, fpt, err := GenerateTradeV2SwapFromTokenToToken(frontRunCtx, w3c, nil, &tf.FrontRunTrade)
	if err != nil {
		log.Warn().Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: failed to generate front run tx")
		log.Err(err).Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: failed to generate front run tx")
		return nil, err
	}
	scInfoFrontRun, err := universalRouterCmdToUnsignedTxPayload(frontRunCtx, w3c, ur)
	if err != nil {
		log.Warn().Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: failed building ur tx")
		log.Err(err).Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: error building ur tx")
		return nil, err
	}
	err = w3c.SuggestAndSetGasPriceAndLimitForTx(ctx, scInfoFrontRun, common.HexToAddress(scInfoFrontRun.SmartContractAddr))
	if err != nil {
		log.Warn().Err(err).Msg("Send: SuggestAndSetGasPriceAndLimitForTx")
		log.Ctx(ctx).Err(err).Msg("Send: SuggestAndSetGasPriceAndLimitForTx")
		return nil, err
	}
	scInfoFrontRun.GasTipCap = artemis_eth_units.NewBigInt(0)
	frontRunTx, err := w3c.GetSignedTxToCallFunctionWithData(ctx, scInfoFrontRun, scInfoFrontRun.Data)
	if err != nil {
		log.Warn().Msg("FRONT_RUN: w3c.GetSignedTxToCallFunctionWithData: error getting signed tx to call function with data")
		log.Err(err).Msg("FRONT_RUN: error getting signed tx to call function with data")
		return nil, err
	}
	frontRunGasCost := artemis_eth_units.MulBigInt(artemis_eth_units.AddBigInt(scInfoFrontRun.GasFeeCap, scInfoFrontRun.GasTipCap), artemis_eth_units.NewBigIntFromUint(scInfoFrontRun.GasLimit))
	tf.FrontRunTrade.TotalGasCost = frontRunGasCost.Uint64()
	log.Info().Uint64("frontRunGasCost", frontRunGasCost.Uint64()).Msg("PackageSandwich: FRONT_RUN gas cost")
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: FRONT_RUN start")

	frTx := TxWithMetadata{
		TradeType: FrontRun,
		ScPayload: scInfoFrontRun,
		Tx:        frontRunTx,
	}
	if fpt != nil {
		frTx.Permit2Tx = fpt.Permit2Tx
	}
	return &frTx, nil
}

func packageBackRun(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow, frScInfo *web3_actions.SendContractTxPayload) (*TxWithMetadata, error) {
	if frScInfo == nil {
		return nil, errors.New("PackageSandwich: BACK_RUN: frScInfo is nil")
	}
	backRunCtx := CreateBackRunCtx(context.Background(), w3c)
	ur, spt, err := GenerateTradeV2SwapFromTokenToToken(backRunCtx, w3c, nil, &tf.SandwichTrade)
	if err != nil {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed to generate sandwich tx")
		log.Err(err).Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed to generate sandwich tx")
		return nil, err
	}
	scInfoSand, err := universalRouterCmdToUnsignedTxPayload(backRunCtx, w3c, ur)
	if err != nil {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed building ur tx")
		log.Err(err).Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed to add tx to bundle group")
		return nil, err
	}
	scInfoSand.GasLimit = frScInfo.GasLimit * 2
	scInfoSand.GasTipCap = artemis_eth_units.MulBigIntFromInt(frScInfo.GasFeeCap, 2)
	scInfoSand.GasFeeCap = artemis_eth_units.MulBigIntFromInt(frScInfo.GasFeeCap, 2)
	signedSandwichTx, err := w3c.GetSignedTxToCallFunctionWithData(backRunCtx, scInfoSand, scInfoSand.Data)
	if err != nil {
		log.Warn().Msg("PackageSandwich: SANDWICH_TRADE: w3c.GetSignedTxToCallFunctionWithData: error getting signed tx to call function with data")
		log.Err(err).Msg("PackageSandwich: SANDWICH_TRADE: error getting signed tx to call function with data")
		return nil, err
	}
	sandwichTx := TxWithMetadata{
		TradeType: BackRun,
		Tx:        signedSandwichTx,
	}
	if spt != nil {
		sandwichTx.Permit2Tx = spt.Permit2Tx
	}
	sandwichGasCost := artemis_eth_units.MulBigInt(artemis_eth_units.AddBigInt(scInfoSand.GasFeeCap, scInfoSand.GasTipCap), artemis_eth_units.NewBigIntFromUint(scInfoSand.GasLimit))
	tf.SandwichTrade.TotalGasCost = sandwichGasCost.Uint64()
	log.Info().Uint64("sandwichGasCost", sandwichGasCost.Uint64()).Msg("PackageSandwich: SANDWICH_TRADE gas cost")
	return &sandwichTx, nil
}

func PackageSandwich(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow) (*MevTxGroup, error) {
	if tf == nil || tf.Tx == nil {
		return nil, errors.New("PackageSandwich: tf is nil")
	}
	if tf.FrontRunTrade.AmountIn == nil || tf.SandwichTrade.AmountOut == nil {
		return nil, errors.New("PackageSandwich: tf.FrontRunTrade.AmountIn or tf.SandwichTrade.AmountOut is nil")
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: start")
	bundle := &MevTxGroup{
		EventID:      0,
		OrderedTxs:   []TxWithMetadata{},
		MevTxs:       []artemis_eth_txs.EthTx{},
		TotalGasCost: artemis_eth_units.NewBigInt(0),
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: FRONT_RUN start")
	frTx, err := packageFrontRun(ctx, w3c, tf)
	if err != nil {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: FRONT_RUN: failed to package front run")
		return nil, err
	}
	if frTx == nil {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: FRONT_RUN: front run tx is nil")
		return nil, errors.New("PackageSandwich: FRONT_RUN: front run tx is nil")
	}
	bundle, err = AddTxToBundleGroup(ctx, *frTx, bundle)
	if err != nil {
		log.Err(err).Interface("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: FRONT_RUN: error adding tx to bundle group")
		return nil, err
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: FRONT_RUN done")

	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: USER_TRADE start")
	// user trade
	userCtx := CreateUserTradeCtx(context.Background())
	userTx := TxWithMetadata{
		TradeType: UserTrade,
		Tx:        tf.Tx,
	}
	bundle, err = AddTxToBundleGroup(userCtx, userTx, bundle)
	if err != nil {
		log.Err(err).Str("txHash", tf.Tx.Hash().String()).Interface("mevTx", bundle.MevTxs).Msg("PackageSandwich: USER_TRADE: failed to add tx to bundle group")
		return nil, err
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: USER_TRADE done")
	// sandwich trade
	if frTx.ScPayload == nil {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: frTx.ScPayload is nil")
		return nil, err
	}
	backRunCtx := CreateBackRunCtx(context.Background(), w3c)
	sandwichTx, err := packageBackRun(ctx, w3c, tf, frTx.ScPayload)
	if err != nil {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed to package back run")
		log.Err(err).Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed to package back run")
		return nil, err
	}
	if sandwichTx == nil {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: sandwich tx is nil")
		return nil, errors.New("PackageSandwich: SANDWICH_TRADE: sandwich tx is nil")
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE start")
	bundle.TotalGasCost = artemis_eth_units.AddBigInt(bundle.TotalGasCost, artemis_eth_units.NewBigIntFromUint(tf.SandwichTrade.TotalGasCost))
	bundle, err = AddTxToBundleGroup(backRunCtx, *sandwichTx, bundle)
	if err != nil {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Interface("mevTx", bundle.MevTxs).Msg("PackageSandwich: SANDWICH_TRADE: failed to add tx to bundle group")
		log.Err(err).Str("txHash", tf.Tx.Hash().String()).Interface("mevTx", bundle.MevTxs).Msg("PackageSandwich: SANDWICH_TRADE: failed to add tx to bundle group")
		return nil, err
	}
	if len(bundle.MevTxs) != 3 {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Int("bundleTxCount", len(bundle.MevTxs)).Msg("PackageSandwich: SANDWICH_TRADE: sandwich bundle not 3 txs")
		return nil, errors.New("PackageSandwich: sandwich bundle not 3 txs")
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE done")
	return bundle, err
}

func StagingPackageSandwichAndCall(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow) (*flashbotsrpc.FlashbotsCallBundleResponse, *MevTxGroup, error) {
	if tf == nil || tf.Tx == nil {
		return nil, nil, errors.New("PackageSandwich: tf is nil")
	}
	if tf.FrontRunTrade.AmountIn == nil || tf.SandwichTrade.AmountOut == nil {
		return nil, nil, errors.New("PackageSandwich: tf.FrontRunTrade.AmountIn or tf.SandwichTrade.AmountOut is nil")
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("StagingPackageSandwichAndCall: start")
	bundle, err := PackageSandwich(ctx, w3c, tf)
	if err != nil {
		log.Err(err).Msg("StagingPackageSandwichAndCall: failed to package sandwich")
		return nil, nil, err
	}
	if bundle == nil {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("StagingPackageSandwichAndCall: bundle is nil")
		return nil, nil, errors.New("bundle is nil")
	}
	//log.Info().Interface("bundle", bundle).Msg("isBundleProfitHigherThanGasFee: bundle")
	//ok, err := isBundleProfitHigherThanGasFee(bundle, tf)
	//if err != nil {
	//	log.Err(err).Bool("ok", ok).Msg("StagingPackageSandwichAndCall: isBundleProfitHigherThanGasFee: failed to check if profit is higher than gas fee")
	//	return nil, nil, err
	//}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("CallFlashbotsBundleStaging: start")
	resp, err := CallFlashbotsBundleStaging(ctx, w3c, *bundle)
	if err != nil {
		log.Err(err).Interface("fbCallResp", resp).Msg("failed to send sandwich")
		return nil, nil, err
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("CallFlashbotsBundleStaging: done")
	log.Info().Str("txHash", tf.Tx.Hash().String()).Interface("fbCallResp", resp).Msg("sent sandwich")
	return &resp, bundle, err
}

func PackageSandwichAndSend(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow) (*flashbotsrpc.FlashbotsSendBundleResponse, error) {
	bundle, err := PackageSandwich(ctx, w3c, tf)
	if err != nil {
		log.Err(err).Msg("failed to package sandwich")
		return nil, err
	}
	if bundle == nil {
		return nil, errors.New("bundle is nil")
	}
	_, err = CallAndSendFlashbotsBundle(ctx, w3c, *bundle)
	if err != nil {
		log.Err(err).Msg("failed to send sandwich")
		return nil, err
	}
	return nil, err
}
