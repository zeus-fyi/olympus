package artemis_flashbots

import (
	"context"

	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
)

func (f *FlashbotsClient) GetBundleStats(ctx context.Context, bundle flashbotsrpc.FlashbotsGetBundleStatsParam) (flashbotsrpc.FlashbotsGetBundleStatsResponse, error) {
	resp, err := f.FlashbotsRPC.FlashbotsGetBundleStats(f.EcdsaPrivateKey(), bundle)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("FlashbotsClient: FlashbotsSendBundle")
		return flashbotsrpc.FlashbotsGetBundleStatsResponse{}, err
	}
	return resp, nil
}
