package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) CreateOrAddToFlashbotsBundle(ur *web3_client.UniversalRouterExecCmd, bn string) error {
	if a.Bundle.FlashbotsSendBundleRequest == nil {
		maxTime := ur.Deadline.Uint64()
		a.Bundle.FlashbotsSendBundleRequest = &flashbotsrpc.FlashbotsSendBundleRequest{
			BlockNumber:  bn,
			MaxTimestamp: &maxTime,
		}
	}
	err := a.Bundle.AddTxs(a.MevTxGroup.GetRawOrderedTxs()...)
	if err != nil {
		return err
	}
	a.trackTxs(a.MevTxGroup)
	a.MevTxGroup.OrderedTxs = []TxWithMetadata{}
	return err
}

func (a *AuxiliaryTradingUtils) CallAndSendFlashbotsBundle(ctx context.Context) (flashbotsrpc.FlashbotsSendBundleResponse, error) {
	sr := flashbotsrpc.FlashbotsSendBundleResponse{}
	eventID, err := a.getBlockNumber(ctx)
	if err != nil {
		log.Err(err).Msg("error getting event id")
		return flashbotsrpc.FlashbotsSendBundleResponse{}, err
	}
	bnStr := hexutil.EncodeUint64(uint64(eventID + 1))
	ctx = a.setBlockNumberCtx(ctx, bnStr)
	_, err = a.callFlashbotsBundle(ctx)
	if err != nil {
		log.Err(err).Msg("error calling flashbots bundle")
		return sr, err
	}
	dbTx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("error beginning db transaction")
		return sr, err
	}
	defer dbTx.Rollback(ctx)
	sr, err = a.sendFlashbotsBundle(ctx)
	if err != nil {
		log.Err(err).Msg("error sending flashbots bundle")
		terr := dbTx.Rollback(ctx)
		if terr != nil {
			log.Err(terr).Msg("error rolling back db transaction")
		}
		return sr, err
	}
	err = artemis_eth_txs.InsertTxsWithBundle(ctx, dbTx, a.MevTxGroup.MevTxs, sr.BundleHash)
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

func (a *AuxiliaryTradingUtils) callFlashbotsBundle(ctx context.Context) (flashbotsrpc.FlashbotsCallBundleResponse, error) {
	var txsCall []string
	bnStr := a.getBlockNumberCtx(ctx)
	txsCall, a.Bundle.Txs = a.Bundle.Txs, txsCall
	callBundle := flashbotsrpc.FlashbotsCallBundleParam{
		Txs:         txsCall,
		BlockNumber: bnStr,
		Timestamp:   a.GetDeadline().Int64(),
	}
	ctx = a.setBlockNumberCtx(ctx, bnStr)
	resp, err := a.f.CallBundle(ctx, callBundle)
	if err != nil {
		log.Err(err).Msg("error calling flashbots bundle")
		a.Bundle.Txs = txsCall
		return resp, err
	}
	a.Bundle.Txs = txsCall
	return resp, nil
}

func (a *AuxiliaryTradingUtils) sendFlashbotsBundle(ctx context.Context) (flashbotsrpc.FlashbotsSendBundleResponse, error) {
	if a.Bundle.FlashbotsSendBundleRequest == nil {
		return flashbotsrpc.FlashbotsSendBundleResponse{}, errors.New("no bundle to send")
	}
	if len(a.Bundle.Txs) == 0 {
		return flashbotsrpc.FlashbotsSendBundleResponse{}, errors.New("no txs to send")
	}
	bundle := &flashbotsrpc.FlashbotsSendBundleRequest{
		Txs:          a.Bundle.Txs,
		BlockNumber:  a.getBlockNumberCtx(ctx),
		MaxTimestamp: a.Bundle.MaxTimestamp,
	}
	resp, err := a.f.SendBundle(ctx, *bundle)
	if err != nil {
		a.Bundle.FlashbotsSendBundleRequest = bundle
		log.Err(err).Msg("error calling flashbots bundle")
		return resp, err
	}
	return resp, nil
}
