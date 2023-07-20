package artemis_flashbots

import (
	"context"
	"errors"

	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
)

func (f *FlashbotsClient) SendBundle(ctx context.Context, bundle flashbotsrpc.FlashbotsSendBundleRequest) (flashbotsrpc.FlashbotsSendBundleResponse, error) {
	resp, err := f.FlashbotsRPC.FlashbotsSendBundle(f.getPrivateKey(), bundle)
	if err != nil {
		log.Warn().Msg("FlashbotsClient: FlashbotsSendBundle")
		log.Ctx(ctx).Error().Err(err).Msg("FlashbotsClient: FlashbotsSendBundle")
		return flashbotsrpc.FlashbotsSendBundleResponse{}, err
	}
	return resp, nil
}

func (f *FlashbotsClient) CallBundle(ctx context.Context, bundle flashbotsrpc.FlashbotsCallBundleParam) (flashbotsrpc.FlashbotsCallBundleResponse, error) {
	if len(bundle.Txs) == 0 {
		return flashbotsrpc.FlashbotsCallBundleResponse{}, errors.New("FlashbotsClient: CallBundle: no txs in bundle")
	}
	if bundle.BlockNumber == "" {
		return flashbotsrpc.FlashbotsCallBundleResponse{}, errors.New("FlashbotsClient: CallBundle: no block number in bundle")
	}
	if bundle.StateBlockNumber == "" {
		bundle.StateBlockNumber = artemis_trading_constants.LatestBlockNumber
	}
	resp, err := f.FlashbotsRPC.FlashbotsCallBundle(f.getPrivateKey(), bundle)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("FlashbotsClient: CallBundle")
		return flashbotsrpc.FlashbotsCallBundleResponse{}, err
	}
	return resp, nil
}
