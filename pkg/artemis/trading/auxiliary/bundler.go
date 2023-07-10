package artemis_trading_auxiliary

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) CreateOrAddToFlashbotsBundle(ur *web3_client.UniversalRouterExecCmd, bn string) {
	if a.Bundle.Txs == nil {
		maxTime := ur.Deadline.Uint64()
		a.Bundle = artemis_flashbots.MevTxBundle{
			FlashbotsSendBundleRequest: flashbotsrpc.FlashbotsSendBundleRequest{
				Txs:          []string{},
				BlockNumber:  bn,
				MaxTimestamp: &maxTime,
			},
		}
	}
	a.Bundle.AddTxs(a.OrderedTxs...)
	a.trackTxs(a.OrderedTxs...)
	a.OrderedTxs = []*types.Transaction{}
}

//  todo, update with missing params

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
		BlockNumber: fmt.Sprintf("%x", eventID),
	}
	resp, err := a.CallBundle(ctx, callBundle)
	if err != nil {
		log.Err(err).Msg("error calling flashbots bundle")
		a.Bundle.Txs = txsCall
		return resp, err
	}
	return resp, nil
}
