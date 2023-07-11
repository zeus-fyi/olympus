package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
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
	err := a.Bundle.AddTxs(a.MevTxGroup.OrderedTxs...)
	if err != nil {
		return err
	}
	a.trackTxs(a.MevTxGroup)
	a.MevTxGroup.OrderedTxs = []*types.Transaction{}
	return err
}

func (a *AuxiliaryTradingUtils) CallAndSendFlashbotsBundle(ctx context.Context) (flashbotsrpc.FlashbotsSendBundleResponse, error) {
	sr := flashbotsrpc.FlashbotsSendBundleResponse{}
	_, err := a.callFlashbotsBundle(ctx)
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
	var txs []artemis_eth_txs.EthTx
	err = artemis_eth_txs.InsertTxsWithBundle(ctx, dbTx, txs, sr.BundleHash)
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
	eventID, err := a.getBlockNumber(ctx)
	if err != nil {
		log.Err(err).Msg("error getting event id")
		return flashbotsrpc.FlashbotsCallBundleResponse{}, err
	}
	txsCall, a.Bundle.Txs = a.Bundle.Txs, txsCall
	callBundle := flashbotsrpc.FlashbotsCallBundleParam{
		Txs:         txsCall,
		BlockNumber: hexutil.EncodeUint64(uint64(eventID + 1)),
		Timestamp:   a.GetDeadline().Int64(),
	}
	resp, err := a.CallBundle(ctx, callBundle)
	if err != nil {
		log.Err(err).Msg("error calling flashbots bundle")
		a.Bundle.Txs = txsCall
		return resp, err
	}
	return resp, nil
}

func (a *AuxiliaryTradingUtils) sendFlashbotsBundle(ctx context.Context) (flashbotsrpc.FlashbotsSendBundleResponse, error) {
	if a.Bundle.FlashbotsSendBundleRequest == nil {
		return flashbotsrpc.FlashbotsSendBundleResponse{}, errors.New("no bundle to send")
	}
	var bundle *flashbotsrpc.FlashbotsSendBundleRequest
	bundle, a.Bundle.FlashbotsSendBundleRequest = a.Bundle.FlashbotsSendBundleRequest, bundle
	resp, err := a.SendBundle(ctx, *bundle)
	if err != nil {
		a.Bundle.FlashbotsSendBundleRequest = bundle
		log.Err(err).Msg("error calling flashbots bundle")
		return resp, err
	}
	return resp, nil
}
