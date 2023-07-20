package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func CallAndSendFlashbotsBundle(ctx context.Context, w3c web3_client.Web3Client, bundle MevTxGroup) (flashbotsrpc.FlashbotsSendBundleResponse, error) {
	sr := flashbotsrpc.FlashbotsSendBundleResponse{}
	eventID, err := getBlockNumber(ctx, w3c)
	if err != nil {
		log.Err(err).Msg("error getting event id")
		return flashbotsrpc.FlashbotsSendBundleResponse{}, err
	}
	bnStr := hexutil.EncodeUint64(uint64(eventID + 1))
	ctx = setBlockNumberCtx(ctx, bnStr)
	_, err = CallFlashbotsBundle(ctx, w3c, &bundle)
	if err != nil {
		log.Err(err).Msg("error calling flashbots bundle")
		return sr, err
	}
	log.Info().Msg("CallFlashbotsBundleStaging: bundle sent successfully")
	dbTx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("error beginning db transaction")
		return sr, err
	}
	defer dbTx.Rollback(ctx)
	sr, err = sendFlashbotsBundle(ctx, w3c, &bundle)
	if err != nil {
		log.Err(err).Msg("error sending flashbots bundle")
		return sr, err
	}
	err = artemis_eth_txs.InsertTxsWithBundle(ctx, dbTx, bundle.MevTxs, sr.BundleHash)
	if err != nil {
		log.Err(err).Msg("error inserting txs with bundle")
		terr := dbTx.Rollback(ctx)
		if terr != nil {
			log.Err(terr).Msg("error rolling back db transaction")
		}
		return sr, err
	}
	err = dbTx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("error committing db transaction")
		return sr, err
	}
	return sr, nil
}

func CallFlashbotsBundleStaging(ctx context.Context, w3c web3_client.Web3Client, bundle MevTxGroup) (flashbotsrpc.FlashbotsCallBundleResponse, error) {
	sr := flashbotsrpc.FlashbotsCallBundleResponse{}
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
		return sr, err
	}
	log.Info().Msg("CallFlashbotsBundleStaging: bundle sent successfully")
	dbTx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("error beginning db transaction")
		return sr, err
	}
	defer dbTx.Rollback(ctx)
	err = artemis_eth_txs.InsertTxsWithBundle(ctx, dbTx, bundle.MevTxs, sr.BundleHash)
	if err != nil {
		log.Info().Str("bundleHash", sr.BundleHash).Interface("bundle.MevTxs", bundle.MevTxs).Msg("CallFlashbotsBundleStaging: error inserting txs with bundle")
		log.Err(err).Msg("error inserting txs with bundle")
		terr := dbTx.Rollback(ctx)
		if terr != nil {
			log.Err(terr).Msg("error rolling back db transaction")
		}
		return sr, err
	}
	err = dbTx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("error committing db transaction")
		return sr, err
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
	f := artemis_flashbots.InitFlashbotsClient(ctx, &w3c.Web3Actions)
	resp, err := f.CallBundle(ctx, fbCallBundle)
	if err != nil {
		log.Warn().Msg("CallFlashbotsBundle: error calling flashbots bundle")
		log.Err(err).Msg("error calling flashbots bundle")
		return resp, err
	}
	return resp, nil
}

func sendFlashbotsBundle(ctx context.Context, w3c web3_client.Web3Client, bundle *MevTxGroup) (flashbotsrpc.FlashbotsSendBundleResponse, error) {
	if bundle == nil || len(bundle.MevTxs) == 0 {
		return flashbotsrpc.FlashbotsSendBundleResponse{}, errors.New("no txs to send or bundle is nil")
	}
	txHexEncodedStrSlice, err := bundle.GetHexEncodedTxStrSlice()
	if err != nil {
		return flashbotsrpc.FlashbotsSendBundleResponse{}, err
	}
	mt := GetDeadline().Uint64()
	fbSendBundle := flashbotsrpc.FlashbotsSendBundleRequest{
		Txs:          txHexEncodedStrSlice,
		BlockNumber:  getBlockNumberCtx(ctx, w3c),
		MaxTimestamp: &mt,
	}
	f := artemis_flashbots.InitFlashbotsClient(ctx, &w3c.Web3Actions)
	resp, err := f.SendBundle(ctx, fbSendBundle)
	if err != nil {
		log.Err(err).Msg("error calling flashbots bundle")
		return resp, err
	}
	return resp, nil
}
