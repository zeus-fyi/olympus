package artemis_flashbots

import (
	"context"

	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
)

/*
TODO: implement setters for the following fields:

	Txs          []string  `json:"txs"`                         // Array[String], A list of signed transactions to execute in an atomic bundle
	BlockNumber  string    `json:"blockNumber"`                 // String, a hex encoded block number for which this bundle is valid on
	MinTimestamp *uint64   `json:"minTimestamp,omitempty"`      // (Optional) Number, the minimum timestamp for which this bundle is valid, in seconds since the unix epoch
	MaxTimestamp *uint64   `json:"maxTimestamp,omitempty"`      // (Optional) Number, the maximum timestamp for which this bundle is valid, in seconds since the unix epoch
	RevertingTxs *[]string `json:"revertingTxHashes,omitempty"` // (Optional) Array[String], A list of tx hashes that are allowed to revert
*/

// {"level":"error","error":"relay error response: unable to decode txs","time":1688709506,"message":"failed to send flashbots bundle"}

func (f *FlashbotsClient) SendBundle(ctx context.Context, bundle flashbotsrpc.FlashbotsSendBundleRequest) (flashbotsrpc.FlashbotsSendBundleResponse, error) {
	resp, err := f.FlashbotsRPC.FlashbotsSendBundle(f.EcdsaPrivateKey(), bundle)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("FlashbotsClient: FlashbotsSendBundle")
		return flashbotsrpc.FlashbotsSendBundleResponse{}, err
	}
	return resp, nil
}

func (f *FlashbotsClient) CallBundle(ctx context.Context, bundle flashbotsrpc.FlashbotsCallBundleParam) (flashbotsrpc.FlashbotsCallBundleResponse, error) {
	resp, err := f.FlashbotsRPC.FlashbotsCallBundle(f.EcdsaPrivateKey(), bundle)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("FlashbotsClient: CallBundle")
		return flashbotsrpc.FlashbotsCallBundleResponse{}, err
	}
	log.Info().Msgf("FlashbotsClient: CallBundle: %v", resp)
	return resp, nil
}

// FlashbotsCallBundle
