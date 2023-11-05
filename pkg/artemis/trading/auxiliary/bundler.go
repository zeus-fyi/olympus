package artemis_trading_auxiliary

import (
	"context"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_reporting "github.com/zeus-fyi/olympus/pkg/artemis/trading/reporting"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func CallAndSendFlashbotsBundle(ctx context.Context, w3c web3_client.Web3Client, bundle MevTxGroup, tf *web3_client.TradeExecutionFlow) (flashbotsrpc.FlashbotsSendBundleResponse, error) {
	eventID, err := getBlockNumber(ctx, w3c)
	if err != nil {
		log.Err(err).Msg("CallAndSendFlashbotsBundle: error getting event id")
		return flashbotsrpc.FlashbotsSendBundleResponse{}, err
	}
	bnStr := hexutil.EncodeUint64(uint64(eventID + 1))
	ctx = setBlockNumberCtx(ctx, bnStr)
	resp, err := CallFlashbotsBundle(ctx, w3c, &bundle)
	if err != nil {
		log.Err(err).Msg("CallAndSendFlashbotsBundle: error calling flashbots bundle")
		return flashbotsrpc.FlashbotsSendBundleResponse{}, err
	}
	log.Info().Int("bn", eventID).Str("bundleHash", resp.BundleHash).Msg("CallAndSendFlashbotsBundle: call bundle simulated successfully")
	gasFees := artemis_eth_units.NewBigIntFromStr(resp.GasFees)
	// 0.005, ~$10
	profitMin := artemis_eth_units.GweiMultiple(5000000)
	profitMarginMin := artemis_eth_units.AddBigInt(gasFees, profitMin)
	log.Info().Interface("resp", resp).Str("tf.SandwichPrediction.ExpectedProfit.String()", tf.SandwichPrediction.ExpectedProfit.String()).
		Str("profitMarginMin", profitMarginMin.String()).Str("gasFees", resp.GasFees).
		Str("resp.BundleGasPrice", resp.BundleGasPrice).Int64("resp.TotalGasUsed", resp.TotalGasUsed).
		Msg("CallAndSendFlashbotsBundle: gas fees vs expected profit")
	if artemis_eth_units.IsXLessThanY(tf.SandwichPrediction.ExpectedProfit, profitMarginMin) {
		log.Warn().Interface("resp", resp).
			Str("expectedProfit", tf.SandwichPrediction.ExpectedProfit.String()).Str("profitMarginMin",
			profitMarginMin.String()).Str("gasFees", resp.GasFees).Str("resp.BundleGasPrice", resp.BundleGasPrice).
			Int64("resp.TotalGasUsed", resp.TotalGasUsed).Msg("CallAndSendFlashbotsBundle: gas fees are greater than expected profit")
		return flashbotsrpc.FlashbotsSendBundleResponse{}, errors.New("CallAndSendFlashbotsBundle: gas fees are greater than expected profit")
	}
	dbTx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("CallAndSendFlashbotsBundle: error beginning db transaction")
		return flashbotsrpc.FlashbotsSendBundleResponse{}, err
	}
	defer dbTx.Rollback(ctx)
	sr, err := sendFlashbotsBundle(ctx, w3c, &bundle)
	if err != nil {
		log.Warn().Msg("CallAndSendFlashbotsBundle: error sending flashbots bundle")
		log.Err(err).Msg("CallAndSendFlashbotsBundle: error sending flashbots bundle")
		return flashbotsrpc.FlashbotsSendBundleResponse{}, err
	}
	log.Info().Int("bn", eventID).Str("bundleHash", sr.BundleHash).Msg("CallAndSendFlashbotsBundle: bundle sent successfully")
	err = artemis_eth_txs.InsertTxsWithBundle(ctx, dbTx, bundle.MevTxs, sr.BundleHash)
	if err != nil {
		log.Err(err).Msg("CallAndSendFlashbotsBundle: error inserting txs with bundle")
		terr := dbTx.Rollback(ctx)
		if terr != nil {
			log.Err(terr).Msg("CallAndSendFlashbotsBundle: error rolling back db transaction")
		}
		return sr, err
	}
	err = dbTx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("CallAndSendFlashbotsBundle: error committing db transaction")
		return sr, err
	}
	return sr, nil
}

func CallFlashbotsBundleStaging(ctx context.Context, w3c web3_client.Web3Client, bundle MevTxGroup) (flashbotsrpc.FlashbotsCallBundleResponse, error) {
	eventID, err := getBlockNumber(ctx, w3c)
	if err != nil {
		log.Warn().Msg("CallFlashbotsBundleStaging: error getting event id")
		log.Err(err).Msg("error getting event id")
		return flashbotsrpc.FlashbotsCallBundleResponse{}, err
	}
	bnStr := hexutil.EncodeUint64(uint64(eventID + 1))
	ctx = setBlockNumberCtx(ctx, bnStr)
	resp, err := CallFlashbotsBundle(ctx, w3c, &bundle)
	if err != nil {
		log.Warn().Msg("CallFlashbotsBundleStaging: error calling flashbots bundle")
		log.Err(err).Msg("error calling flashbots bundle")
		return resp, err
	}
	log.Info().Int("bn", eventID).Str("bundleHash", resp.BundleHash).Msg("CallFlashbotsBundleStaging: bundle sent successfully")
	dbTx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("error beginning db transaction")
		return resp, err
	}
	defer dbTx.Rollback(ctx)
	err = artemis_eth_txs.InsertTxsWithBundle(ctx, dbTx, bundle.MevTxs, resp.BundleHash)
	if err != nil {
		log.Info().Str("bundleHash", resp.BundleHash).Interface("bundle.MevTxs", bundle.MevTxs).Msg("CallFlashbotsBundleStaging: error inserting txs with bundle")
		log.Err(err).Msg("error inserting txs with bundle")
		terr := dbTx.Rollback(ctx)
		if terr != nil {
			log.Err(terr).Msg("error rolling back db transaction")
		}
		return resp, err
	}
	err = dbTx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("error committing db transaction")
		return resp, err
	}
	return resp, nil
}

func CallFlashbotsBundle(ctx context.Context, w3c web3_client.Web3Client, bundle *MevTxGroup) (flashbotsrpc.FlashbotsCallBundleResponse, error) {
	if bundle == nil || len(bundle.MevTxs) == 0 {
		return flashbotsrpc.FlashbotsCallBundleResponse{}, errors.New("no txs to send or bundle is nil")
	}
	bnStr := getBlockNumberCtx(ctx, w3c)
	txHexEncodedStrSlice, err := bundle.GetHexEncodedTxStrSlice()
	if err != nil {
		log.Warn().Msg("CallFlashbotsBundle: error getting hex encoded tx str slice")
		return flashbotsrpc.FlashbotsCallBundleResponse{}, err
	}
	fbCallBundle := flashbotsrpc.FlashbotsCallBundleParam{
		Txs:         txHexEncodedStrSlice,
		BlockNumber: bnStr,
		Timestamp:   GetDeadline().Int64(),
	}
	ctx = setBlockNumberCtx(ctx, bnStr)
	sendAdditionalCallBundles(ctx, w3c, fbCallBundle)
	f := artemis_flashbots.InitFlashbotsClient(ctx, &w3c.Web3Actions)
	resp, err := f.CallBundle(ctx, fbCallBundle)
	if err != nil {
		log.Warn().Msg("CallFlashbotsBundle: error calling flashbots bundle")
		log.Err(err).Msg("error calling flashbots bundle")
		return resp, err
	}
	go func() {
		err = artemis_reporting.InsertCallBundleResp(ctx, "flashbots", 1, resp)
		if err != nil {
			log.Warn().Msg("CallFlashbotsBundle: error inserting call bundle resp")
			log.Err(err).Msg("error inserting call bundle resp")
			return
		}
	}()

	log.Info().Interface("resp", resp).Str("resp.BundleGasPrice", resp.BundleGasPrice).Interface("fbCallResp", resp.Results).Msg("CallFlashbotsBundle: bundle sent successfully")
	return resp, nil
}

func sendFlashbotsBundle(ctx context.Context, w3c web3_client.Web3Client, bundle *MevTxGroup) (flashbotsrpc.FlashbotsSendBundleResponse, error) {
	if bundle == nil || len(bundle.MevTxs) == 0 {
		return flashbotsrpc.FlashbotsSendBundleResponse{}, errors.New("no txs to send or bundle is nil")
	}
	txHexEncodedStrSlice, err := bundle.GetHexEncodedTxStrSlice()
	if err != nil {
		log.Err(err).Msg("sendFlashbotsBundle: error getting hex encoded tx str slice")
		return flashbotsrpc.FlashbotsSendBundleResponse{}, err
	}
	mt := GetDeadline().Uint64()
	fbSendBundle := flashbotsrpc.FlashbotsSendBundleRequest{
		Txs:          txHexEncodedStrSlice,
		BlockNumber:  getBlockNumberCtx(ctx, w3c),
		MaxTimestamp: &mt,
	}
	f := artemis_flashbots.InitFlashbotsClient(ctx, &w3c.Web3Actions)
	sendAdditionalBundles(ctx, w3c, fbSendBundle)
	resp, err := f.SendBundle(ctx, fbSendBundle)
	if err != nil {
		log.Warn().Msg("sendFlashbotsBundle: error sending flashbots bundle")
		log.Err(err).Msg("sendFlashbotsBundle: error sending flashbots bundle")
		return resp, err
	}
	log.Info().Str("bundleHash", resp.BundleHash).Msg("sendFlashbotsBundle: bundle sent successfully")
	return resp, nil
}

func sendAdditionalBundles(ctx context.Context, w3c web3_client.Web3Client, fbSendBundle flashbotsrpc.FlashbotsSendBundleRequest) {
	builders := artemis_flashbots.Builders
	for _, builder := range builders {
		f := artemis_flashbots.InitFlashbotsClientForAdditionalBuilder(ctx, &w3c.Web3Actions, builder)
		go func(builder string, f artemis_flashbots.FlashbotsClient) {
			log.Info().Str("builder", builder).Msg("sendAdditionalBundles: sending bundle")
			resp, err := f.SendBundle(ctx, fbSendBundle)
			if err != nil {
				log.Warn().Str("builder", builder).Msg("sendAdditionalBundles: error calling sending bundle")
				log.Err(err).Str("builder", builder).Msg("sendAdditionalBundles: error calling sending bundle")
			}
			log.Info().Str("builder", builder).Str("bundleHash", resp.BundleHash).Msg("sendAdditionalBundles: bundle sent successfully")
		}(builder, f)
	}
}

func sendAdditionalCallBundles(ctx context.Context, w3c web3_client.Web3Client, callBundle flashbotsrpc.FlashbotsCallBundleParam) {
	builders := artemis_flashbots.Builders
	for _, builder := range builders {
		if builder == artemis_flashbots.Builder69 {
			continue
		}
		if builder == artemis_flashbots.BeaverRelay {
			continue
		}
		if builder == artemis_flashbots.PayloadBuilder {
			continue
		}
		if builder == artemis_flashbots.TitanBuilder {
			continue
		}
		if builder == artemis_flashbots.EdenBuilder {
			continue
		}
		if builder == artemis_flashbots.RsyncBuilder {
			continue
		}
		if builder == artemis_flashbots.BuildAIBuilder {
			continue
		}
		if builder == artemis_flashbots.EthBuilder {
			continue
		}
		f := artemis_flashbots.InitFlashbotsClientForAdditionalBuilder(ctx, &w3c.Web3Actions, builder)
		go func(builder string, f artemis_flashbots.FlashbotsClient) {
			log.Info().Str("builder", builder).Msg("sendAdditionalCallBundles: sending bundle")
			resp, err := f.CallBundle(ctx, callBundle)
			if err != nil {
				log.Warn().Str("builder", builder).Msg("sendAdditionalCallBundles: error calling sending bundle")
				log.Err(err).Str("builder", builder).Msg("sendAdditionalCallBundles: error calling sending bundle")
				return
			}

			if strings.Contains(builder, "blocknative") {
				err = artemis_reporting.InsertCallBundleResp(ctx, builder, 1, resp)
				if err != nil {
					log.Warn().Msg("sendAdditionalCallBundles: error inserting call bundle resp")
					log.Err(err).Msg("sendAdditionalCallBundles: error inserting call bundle resp")
					return
				}
			}

			log.Info().Str("builder", builder).Str("bundleHash", resp.BundleHash).Msg("sendAdditionalBundles: bundle sent successfully")
		}(builder, f)
	}
}
