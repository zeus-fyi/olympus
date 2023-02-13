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
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
)

var (
	AttestationSigningRequestPriorityQueue = SignaturePriorityQueue{
		Type:          ATTESTATION,
		PriorityQueue: datastructures.PriorityQueue{},
	}
	AggregationSlotSigningRequestPriorityQueue = SignaturePriorityQueue{
		Type:          AGGREGATION_SLOT,
		PriorityQueue: datastructures.PriorityQueue{},
	}
	AggregationAndProofSigningRequestPriorityQueue = SignaturePriorityQueue{
		Type:          AGGREGATE_AND_PROOF,
		PriorityQueue: datastructures.PriorityQueue{},
	}
	BlockSigningRequestPriorityQueue = SignaturePriorityQueue{
		Type:          BLOCK_V2,
		PriorityQueue: datastructures.PriorityQueue{},
	}
	RandaoRevealSigningRequestPriorityQueue = SignaturePriorityQueue{
		Type:          RANDAO_REVEAL,
		PriorityQueue: datastructures.PriorityQueue{},
	}
	SyncCommitteeMessageSigningRequestPriorityQueue = SignaturePriorityQueue{
		Type:          SYNC_COMMITTEE_MESSAGE,
		PriorityQueue: datastructures.PriorityQueue{},
	}
	SyncCommitteeSelectionProofSigningRequestPriorityQueue = SignaturePriorityQueue{
		Type:          SYNC_COMMITTEE_SELECTION_PROOF,
		PriorityQueue: datastructures.PriorityQueue{},
	}
	SyncCommitteeContributionAndProof = SignaturePriorityQueue{
		Type:          SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF,
		PriorityQueue: datastructures.PriorityQueue{},
	}
	ValidatorRegistration = SignaturePriorityQueue{
		Type:          VALIDATOR_REGISTRATION,
		PriorityQueue: datastructures.PriorityQueue{},
	}

	SignatureResponsesCache = cache.New(10*time.Second, 20*time.Second)
)

type SignaturePriorityQueue struct {
	Type string
	datastructures.PriorityQueue
}

func InitAsyncMessageQueues(ctx context.Context) {
	for {
		go AttestationSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		time.Sleep(2 * time.Millisecond)
		go AggregationSlotSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		time.Sleep(2 * time.Millisecond)
		go AggregationAndProofSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		time.Sleep(2 * time.Millisecond)
		go BlockSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		time.Sleep(2 * time.Millisecond)
		go RandaoRevealSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		time.Sleep(2 * time.Millisecond)
		go SyncCommitteeMessageSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		time.Sleep(2 * time.Millisecond)
		go SyncCommitteeSelectionProofSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		time.Sleep(2 * time.Millisecond)
		go SyncCommitteeContributionAndProof.SendSignatureRequestsFromQueue(ctx)
		time.Sleep(2 * time.Millisecond)
		go ValidatorRegistration.SendSignatureRequestsFromQueue(ctx)
		time.Sleep(2 * time.Millisecond)
	}
}

func (sq *SignaturePriorityQueue) SendSignatureRequestsFromQueue(ctx context.Context) {
	batchSigReqs := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}
	m := make(map[string]string)
	seen := make(map[string]SignRequest)
	ql := sq.Len()
	if ql == 0 {
		return
	}
	log.Info().Str("signingType", sq.Type).Msg(fmt.Sprintf("queue length: %d", ql))
	for i := 0; i < ql; i++ {
		sr := sq.Pop().(SignRequest)
		pubkey := sr.Pubkey
		if v, ok := seen[pubkey]; ok {
			log.Ctx(ctx).Warn().Interface("prevSignRequest", v).Interface("currentSignRequest", sr).Msg(fmt.Sprintf("more than one message seen for pubkey %s, adding back to the queue", pubkey))
			sq.Push(sr)
		}
		seen[pubkey] = sr
		batchSigReqs.Map[sr.Pubkey] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: sr.SigningRoot}
		m[sr.Pubkey] = sr.UUID.String()
	}
	var resp aegis_inmemdbs.EthereumBLSKeySignatureResponses
	var err error
	switch ethereum_slashing_protection_watermarking.Network {
	case "mainnet":
		if sq.Type == ATTESTATION || sq.Type == AGGREGATION_SLOT || sq.Type == AGGREGATE_AND_PROOF {
			resp, err = eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsMainnetWorker.ExecuteValidatorSignatureRequestsWorkflow(ctx, batchSigReqs)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("ExecuteValidatorSignatureRequestsWorkflow")
			}
		} else {
			resp, err = eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsMainnetWorkerSecondary.ExecuteValidatorSignatureRequestsWorkflow(ctx, batchSigReqs)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("ExecuteValidatorSignatureRequestsWorkflow")
			}
		}
	case "ephemery":
		if sq.Type == ATTESTATION || sq.Type == AGGREGATION_SLOT || sq.Type == AGGREGATE_AND_PROOF {
			resp, err = eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsEphemeryWorker.ExecuteValidatorSignatureRequestsWorkflow(ctx, batchSigReqs)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("ExecuteValidatorSignatureRequestsWorkflow")
			}
		} else {
			resp, err = eth_validator_signature_requests.ArtemisEthereumValidatorSignatureRequestsEphemeryWorkerSecondary.ExecuteValidatorSignatureRequestsWorkflow(ctx, batchSigReqs)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("ExecuteValidatorSignatureRequestsWorkflow")
			}
		}
	}
	for pubkey, msg := range resp.Map {
		sigResp := msg.Signature
		uuid := m[pubkey]
		SignatureResponsesCache.Set(uuid, sigResp, cache.DefaultExpiration)
	}
}

func WaitForSignature(ctx context.Context, sr SignRequest) (Eth2SignResponse, error) {
	ch := make(chan Eth2SignResponse)
	go func(ctx context.Context, sr SignRequest) {
		ch <- ReturnSignedMessage(ctx, sr)
	}(ctx, sr)
	resp := <-ch
	return resp, nil
}

func ReturnSignedMessage(ctx context.Context, sr SignRequest) Eth2SignResponse {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	for {
		if v, found := SignatureResponsesCache.Get(sr.UUID.String()); found {
			sigResp := v.(string)
			return Eth2SignResponse{Signature: sigResp}
		}
		time.Sleep(5 * time.Millisecond)
	}
}
