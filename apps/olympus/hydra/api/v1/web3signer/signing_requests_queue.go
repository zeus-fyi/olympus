package hydra_eth2_web3signer

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	ethereum_slashing_protection_watermarking "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/slashing_protection"
	eth_validator_signature_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests"
	"github.com/zeus-fyi/olympus/pkg/utils/datastructures"
	bls_serverless_signing "github.com/zeus-fyi/zeus/pkg/aegis/aws/serverless_signing"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
)

var (
	SigningRequestPriorityQueue       = datastructures.PriorityQueue{}
	SignatureResponseAggregationCache = cache.New(1*time.Second, 2*time.Second)
)

func QueueSigningRequestAndReturnSignature(ctx context.Context, sr SignRequest) (Eth2SignResponse, error) {
	ch := make(chan Eth2SignResponse)
	SigningRequestPriorityQueue.Push(sr)
	go func(sr SignRequest) {
		ch <- ReturnSignedMessage(ctx, sr)
	}(sr)
	resp := <-ch
	return resp, nil
}

func ReturnSignedMessage(ctx context.Context, sr SignRequest) Eth2SignResponse {

	return Eth2SignResponse{}
}

func ProcessSigningRequestQueue(ctx context.Context) error {
	batchSigReq := bls_serverless_signing.SignatureRequests{
		SignatureRequests: aegis_inmemdbs.EthereumBLSKeySignatureRequests{},
	}
	for i := 0; i < SigningRequestPriorityQueue.Len(); i++ {
		sr := SigningRequestPriorityQueue.Pop().(SignRequest)
		batchSigReq.SignatureRequests.Map[sr.Pubkey] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: sr.SigningRoot}
		log.Ctx(ctx).Info().Interface("sr", sr).Msg("ProcessSigningRequestQueue")
	}

	// TODO, send batchSigReq to serverless signing
	resp := aegis_inmemdbs.EthereumBLSKeySignatureResponses{}
	var err error
	switch ethereum_slashing_protection_watermarking.Network {
	case "mainnet":
		resp, err = eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsMainnetWorker.ExecuteValidatorSignatureRequestsWorkflow(ctx, batchSigReq)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("ExecuteValidatorSignatureRequestsWorkflow")
		}
	case "ephemery":
		resp, err = eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsEphemeryWorker.ExecuteValidatorSignatureRequestsWorkflow(ctx, batchSigReq)
		log.Ctx(ctx).Error().Err(err).Msg("ExecuteValidatorSignatureRequestsWorkflow")
	}

	fmt.Print(resp)

	SignatureResponseAggregationCache.Set("batchSigReq", batchSigReq, cache.DefaultExpiration)
	return nil
}
