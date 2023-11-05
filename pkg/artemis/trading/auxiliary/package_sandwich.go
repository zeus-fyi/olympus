package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func packageFrontRun(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow) (*TxWithMetadata, error) {
	frontRunCtx := CreateFrontRunCtx(context.Background())
	if tf.InitialPair != nil {
		ur, fpt, err := GenerateTradeV2SwapFromTokenToToken(frontRunCtx, w3c, nil, &tf.FrontRunTrade)
		if err != nil {
			log.Warn().Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: failed to generate front run tx")
			log.Err(err).Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: failed to generate front run tx")
			return nil, err
		}
		scInfoFrontRun, err := universalRouterCmdToUnsignedTxPayload(frontRunCtx, ur)
		if err != nil {
			log.Warn().Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: failed building ur tx")
			log.Err(err).Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: error building ur tx")
			return nil, err
		}
		if scInfoFrontRun.Data == nil || len(scInfoFrontRun.Data) <= 0 {
			log.Warn().Msg("FRONT_RUN: scInfoFrontRun.Data is nil")
			return nil, errors.New("FRONT_RUN: scInfoFrontRun.Data is nil")
		}
		err = w3c.SuggestAndSetGasPriceAndLimitForTx(ctx, scInfoFrontRun, common.HexToAddress(scInfoFrontRun.SmartContractAddr))
		if err != nil {
			log.Warn().Err(err).Msg("Send: SuggestAndSetGasPriceAndLimitForTx")
			log.Ctx(ctx).Err(err).Msg("Send: SuggestAndSetGasPriceAndLimitForTx")
			return nil, err
		}
		scInfoFrontRun.GasTipCap = artemis_eth_units.NewBigInt(0)
		scInfoFrontRun.GasFeeCap = artemis_eth_units.MulBigIntWithFloat(scInfoFrontRun.GasFeeCap, 0.75)
		scInfoFrontRun.GasPrice = scInfoFrontRun.GasFeeCap
		scInfoFrontRun.GasLimit = uint64(float64(scInfoFrontRun.GasLimit) * 1.3)
		frontRunTx, err := w3c.GetSignedTxToCallFunctionWithData(ctx, scInfoFrontRun, scInfoFrontRun.Data)
		if err != nil {
			log.Warn().Msg("FRONT_RUN: w3c.GetSignedTxToCallFunctionWithData: error getting signed tx to call function with data")
			log.Err(err).Msg("FRONT_RUN: error getting signed tx to call function with data")
			return nil, err
		}
		log.Info().Str("txHash", tf.Tx.Hash().String()).Str("gasFeeGap", scInfoFrontRun.GasFeeCap.String()).Str("gasTipCap", scInfoFrontRun.GasTipCap.String()).Uint64("gasLimit", scInfoFrontRun.GasLimit).Msg("PackageSandwich: FRONT_RUN gas")
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

	if tf.InitialPairV3 != nil && tf.FrontRunTrade.AmountFees != nil {
		ur, fpt, err := GenerateTradeV3SwapFromTokenToToken(frontRunCtx, w3c, nil, &tf.FrontRunTrade)
		if err != nil {
			log.Warn().Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: failed to generate front run tx")
			log.Err(err).Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: failed to generate front run tx")
			return nil, err
		}
		scInfoFrontRun, err := universalRouterCmdToUnsignedTxPayload(frontRunCtx, ur)
		if err != nil {
			log.Warn().Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: failed building ur tx")
			log.Err(err).Interface("txHash", tf.Tx.Hash().String()).Msg("FRONT_RUN: error building ur tx")
			return nil, err
		}
		if scInfoFrontRun.Data == nil || len(scInfoFrontRun.Data) <= 0 {
			log.Warn().Msg("FRONT_RUN: scInfoFrontRun.Data is nil")
			return nil, errors.New("FRONT_RUN: scInfoFrontRun.Data is nil")
		}
		err = w3c.SuggestAndSetGasPriceAndLimitForTx(ctx, scInfoFrontRun, common.HexToAddress(scInfoFrontRun.SmartContractAddr))
		if err != nil {
			log.Warn().Err(err).Msg("Send: SuggestAndSetGasPriceAndLimitForTx")
			log.Ctx(ctx).Err(err).Msg("Send: SuggestAndSetGasPriceAndLimitForTx")
			return nil, err
		}
		scInfoFrontRun.GasTipCap = artemis_eth_units.NewBigInt(0)
		scInfoFrontRun.GasFeeCap = artemis_eth_units.MulBigIntWithFloat(scInfoFrontRun.GasFeeCap, 0.7)
		scInfoFrontRun.GasPrice = scInfoFrontRun.GasFeeCap
		scInfoFrontRun.GasLimit = uint64(float64(scInfoFrontRun.GasLimit) * 1.2)
		toAddr := common.HexToAddress(scInfoFrontRun.ToAddress.Hex())
		msg := ethereum.CallMsg{
			From:      common.HexToAddress(w3c.Address().Hex()),
			To:        &toAddr,
			GasFeeCap: scInfoFrontRun.GasFeeCap,
			GasTipCap: scInfoFrontRun.GasTipCap,
			Data:      scInfoFrontRun.Data,
			Value:     scInfoFrontRun.Amount,
		}
		gasLimit, err := w3c.C.EstimateGas(ctx, msg)
		if err != nil {
			log.Warn().Err(err).Msg("FRONT_RUN: SuggestAndSetGasPriceAndLimitForTx: EstimateGas")
			log.Ctx(ctx).Err(err).Msg("FRONT_RUN: SuggestAndSetGasPriceAndLimitForTx: EstimateGas")
			return nil, err
		}
		scInfoFrontRun.GasPrice = scInfoFrontRun.GasFeeCap
		scInfoFrontRun.GasLimit = uint64(float64(gasLimit) * 1.3)
		frontRunTx, err := w3c.GetSignedTxToCallFunctionWithData(ctx, scInfoFrontRun, scInfoFrontRun.Data)
		if err != nil {
			log.Warn().Msg("FRONT_RUN: w3c.GetSignedTxToCallFunctionWithData: error getting signed tx to call function with data")
			log.Err(err).Msg("FRONT_RUN: error getting signed tx to call function with data")
			return nil, err
		}
		log.Info().Str("txHash", tf.Tx.Hash().String()).Str("gasFeeGap", scInfoFrontRun.GasFeeCap.String()).Str("gasTipCap", scInfoFrontRun.GasTipCap.String()).Uint64("gasLimit", gasLimit).Msg("PackageSandwich: FRONT_RUN gas ")
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
	return nil, errors.New("PackageSandwich: FRONT_RUN: tf.InitialPair and tf.InitialPairV3 are nil")

}

func packageBackRun(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow, frScInfo *web3_actions.SendContractTxPayload) (*TxWithMetadata, error) {
	if frScInfo == nil {
		return nil, errors.New("PackageSandwich: BACK_RUN: frScInfo is nil")
	}
	if tf.InitialPair != nil {
		ur, spt, err := GenerateTradeV2SwapFromTokenToToken(ctx, w3c, nil, &tf.SandwichTrade)
		if err != nil {
			log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed to generate sandwich tx")
			log.Err(err).Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed to generate sandwich tx")
			return nil, err
		}
		scInfoSand, err := universalRouterCmdToUnsignedTxPayload(ctx, ur)
		if err != nil {
			log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed building ur tx")
			log.Err(err).Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed to add tx to bundle group")
			return nil, err
		}

		scInfoSand.GasLimit = uint64(float64(frScInfo.GasLimit) * 1.1)
		scInfoSand.GasTipCap = artemis_eth_units.MulBigIntWithFloat(frScInfo.GasFeeCap, 1.5)
		scInfoSand.GasFeeCap = scInfoSand.GasTipCap
		scInfoSand.GasPrice = scInfoSand.GasFeeCap
		backRunCtx := CreateBackRunCtx(context.Background())
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
		log.Info().Str("baseFee", frScInfo.GasPrice.String()).Str("gasFeeGap", scInfoSand.GasFeeCap.String()).Str("gasTipCap", scInfoSand.GasTipCap.String()).Uint64("gasLimit", scInfoSand.GasLimit).Msg("PackageSandwich: SANDWICH_TRADE gas")
		return &sandwichTx, nil
	}

	if tf.InitialPairV3 != nil && tf.SandwichTrade.AmountFees != nil {
		ur, spt, err := GenerateTradeV3SwapFromTokenToToken(ctx, w3c, nil, &tf.SandwichTrade)
		if err != nil {
			log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed to generate sandwich tx")
			log.Err(err).Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed to generate sandwich tx")
			return nil, err
		}
		scInfoSand, err := universalRouterCmdToUnsignedTxPayload(ctx, ur)
		if err != nil {
			log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed building ur tx")
			log.Err(err).Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: failed to add tx to bundle group")
			return nil, err
		}
		scInfoSand.GasLimit = uint64(float64(frScInfo.GasLimit) * 1.1)
		scInfoSand.GasTipCap = artemis_eth_units.MulBigIntWithFloat(frScInfo.GasFeeCap, 1.8)
		scInfoSand.GasFeeCap = scInfoSand.GasTipCap
		scInfoSand.GasPrice = scInfoSand.GasTipCap
		backRunCtx := CreateBackRunCtx(context.Background())
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
		log.Info().Str("txHash", tf.Tx.Hash().String()).Str("baseFee", frScInfo.GasPrice.String()).Str("gasFeeGap", scInfoSand.GasFeeCap.String()).Str("gasTipCap", scInfoSand.GasTipCap.String()).Uint64("gasLimit", scInfoSand.GasLimit).Msg("PackageSandwich: SANDWICH_TRADE gas")
		return &sandwichTx, nil
	}
	return nil, errors.New("PackageSandwich: SANDWICH_TRADE: tf.InitialPair and tf.InitialPairV3 are nil")
}

func PackageSandwich(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow) (*MevTxGroup, error) {
	if tf == nil || tf.Tx == nil {
		return nil, errors.New("PackageSandwich: tf is nil")
	}
	if w3c.Account == nil {
		return nil, errors.New("PackageSandwich: w3c.Account is nil")
	}
	if tf.FrontRunTrade.AmountIn == nil || tf.SandwichTrade.AmountOut == nil {
		return nil, errors.New("PackageSandwich: tf.FrontRunTrade.AmountIn or tf.SandwichTrade.AmountOut is nil")
	}

	log.Info().Str("txHash", tf.Tx.Hash().String()).Str("traderAccount", w3c.Account.PublicKey()).Msg("PackageSandwich: start")
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
	log.Info().Uint64("PackageSandwich: FRONT_RUN: tradersNonceFrontRun", frTx.Tx.Nonce())
	bundle := &MevTxGroup{
		EventID:      0,
		OrderedTxs:   []TxWithMetadata{},
		MevTxs:       []artemis_eth_txs.EthTx{},
		TotalGasCost: artemis_eth_units.NewBigInt(0),
		BaseFee:      frTx.ScPayload.GasFeeCap,
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
		BaseFee:   bundle.BaseFee,
	}
	bundle, err = AddTxToBundleGroup(userCtx, userTx, bundle)
	if err != nil {
		log.Err(err).Str("txHash", tf.Tx.Hash().String()).Interface("mevTx", bundle.MevTxs).Msg("PackageSandwich: USER_TRADE: failed to add tx to bundle group")
		return nil, err
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: USER_TRADE done")
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE start")
	// sandwich trade
	if frTx.ScPayload == nil {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwich: SANDWICH_TRADE: frTx.ScPayload is nil")
		return nil, err
	}
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
	log.Info().Uint64("PackageSandwich: SANDWICH_TRADE: tradersNonceBackRun", sandwichTx.Tx.Nonce())
	bundle.TotalGasCost = artemis_eth_units.AddBigInt(bundle.TotalGasCost, artemis_eth_units.NewBigIntFromUint(tf.SandwichTrade.TotalGasCost))
	bundle, err = AddTxToBundleGroup(ctx, *sandwichTx, bundle)
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

func PackageSandwichAndSend(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow, m *metrics_trading.TradingMetrics) (*flashbotsrpc.FlashbotsSendBundleResponse, error) {
	if tf == nil || tf.Tx == nil {
		return nil, errors.New("PackageSandwichAndSend: tf is nil")
	}
	if artemis_eth_units.IsXLessThanEqZeroOrOne(tf.FrontRunTrade.AmountIn) || artemis_eth_units.IsXLessThanEqZeroOrOne(tf.SandwichTrade.AmountOut) {
		return nil, errors.New("PackageSandwichAndSend: tf.FrontRunTrade.AmountIn or tf.SandwichTrade.AmountOut is nil")
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwichAndSend: start")
	bundle, err := PackageSandwich(ctx, w3c, tf)
	if err != nil {
		log.Err(err).Msg("PackageSandwichAndSend: PackageSandwich failed to package sandwich")
		return nil, err
	}
	if bundle == nil {
		return nil, errors.New("bundle is nil")
	}
	resp, err := CallAndSendFlashbotsBundle(ctx, w3c, *bundle, tf)
	if err != nil {
		log.Err(err).Msg("PackageSandwichAndSend: CallAndSendFlashbotsBundle failed to send sandwich")
		return nil, err
	}
	log.Info().Str("bundleHash", resp.BundleHash).Msg("PackageSandwichAndSend: done")
	return &resp, err
}

func ReadOnlyPackageSandwichAndCall(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow, m *metrics_trading.TradingMetrics) error {
	if tf == nil || tf.Tx == nil {
		return errors.New("PackageSandwichAndSend: tf is nil")
	}
	if artemis_eth_units.IsXLessThanEqZeroOrOne(tf.FrontRunTrade.AmountIn) || artemis_eth_units.IsXLessThanEqZeroOrOne(tf.SandwichTrade.AmountOut) {
		return errors.New("PackageSandwichAndSend: tf.FrontRunTrade.AmountIn or tf.SandwichTrade.AmountOut is nil")
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Msg("PackageSandwichAndSend: start")
	bundle, err := PackageSandwich(ctx, w3c, tf)
	if err != nil {
		log.Err(err).Msg("PackageSandwichAndSend: PackageSandwich failed to package sandwich")
		return err
	}
	if bundle == nil {
		return errors.New("bundle is nil")
	}
	err = CallReadOnlyBundle(ctx, w3c, *bundle, tf)
	if err != nil {
		log.Err(err).Msg("PackageSandwichAndSend: CallAndSendFlashbotsBundle failed to send sandwich")
		return err
	}
	return nil
}
