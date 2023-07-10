package artemis_trading_auxiliary

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
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
	err := a.Bundle.AddTxs(a.OrderedTxs...)
	if err != nil {
		return err
	}
	a.trackTxs(a.OrderedTxs...)
	a.OrderedTxs = []*types.Transaction{}
	return err
}

func (a *AuxiliaryTradingUtils) callFlashbotsBundle(ctx context.Context) (flashbotsrpc.FlashbotsCallBundleResponse, error) {
	var txsCall []string
	eventID, err := a.getEventID(ctx)
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
