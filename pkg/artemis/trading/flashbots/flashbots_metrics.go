package artemis_flashbots

import (
	"context"

	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
)

func (f *FlashbotsClient) GetBundleStatsV2(ctx context.Context, bundle flashbotsrpc.FlashbotsGetBundleStatsParam) (flashbotsrpc.FlashbotsGetBundleStatsResponseV2, error) {
	resp, err := f.FlashbotsRPC.FlashbotsGetBundleStatsV2(f.getPrivateKey(), bundle)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("FlashbotsClient: FlashbotsSendBundle")
		return flashbotsrpc.FlashbotsGetBundleStatsResponseV2{}, err
	}
	return resp, nil
}
