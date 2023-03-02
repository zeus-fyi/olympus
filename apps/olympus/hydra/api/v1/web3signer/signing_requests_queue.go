package hydra_eth2_web3signer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
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
	SyncCommitteeContributionAndProofPriorityQueue = SignaturePriorityQueue{
		Type:          SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF,
		PriorityQueue: datastructures.PriorityQueue{},
	}
	ValidatorRegistrationPriorityQueue = SignaturePriorityQueue{
		Type:          VALIDATOR_REGISTRATION,
		PriorityQueue: datastructures.PriorityQueue{},
	}

	SignatureResponsesCache = cache.New(10*time.Second, 20*time.Second)
)

type SignaturePriorityQueue struct {
	Type string
	datastructures.PriorityQueue
}

func InitAsyncMessageQueuesSyncCommitteeQueues(ctx context.Context) {
	minDuration := 30 * time.Millisecond
	maxDuration := 50 * time.Millisecond
	for {
		go func() {
			SyncCommitteeMessageSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		}()
		go func() {
			SyncCommitteeSelectionProofSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		}()
		go func() {
			SyncCommitteeContributionAndProofPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		}()
		jitterDelay := time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration
		time.Sleep(jitterDelay)
	}
}

func InitAsyncMessageAttestationQueues(ctx context.Context) {
	minDuration := 10 * time.Millisecond
	maxDuration := 20 * time.Millisecond
	go func() {
		AttestationSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
	}()
	go func() {
		AggregationSlotSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
	}()
	go func() {
		AggregationAndProofSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
	}()
	jitterDelay := time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration
	time.Sleep(jitterDelay)
}

func InitAsyncBlockMessageQueues(ctx context.Context) {
	minDuration := 10 * time.Millisecond
	maxDuration := 30 * time.Millisecond
	for {
		go func() {
			BlockSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		}()
		go func() {
			RandaoRevealSigningRequestPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		}()
		go func() {
			ValidatorRegistrationPriorityQueue.SendSignatureRequestsFromQueue(ctx)
		}()
		jitterDelay := time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration
		time.Sleep(jitterDelay)
	}
}

func (sq *SignaturePriorityQueue) SendSignatureRequestsFromQueue(ctx context.Context) {
	batchSigReqs := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}
	pubkeyToUUID := make(map[string]string)
	seen := make(map[string]SignRequest)
	ql := sq.Len()
	if ql == 0 {
		return
	}
	log.Info().Str("signingType", sq.Type).Msg(fmt.Sprintf("queue length: %d", ql))
	for i := 0; i < ql; i++ {
		item := sq.Pop().(*datastructures.Item)
		sr := item.Value.(SignRequest)
		pubkey := sr.Pubkey
		if v, ok := seen[pubkey]; ok {
			log.Ctx(ctx).Warn().Interface("prevSignRequest", v).Interface("currentSignRequest", sr).Msg(fmt.Sprintf("more than one message seen for pubkey %s, adding back to the queue", pubkey))
			sq.Push(SigningRequestToItem(sr))
		}
		seen[pubkey] = sr
		batchSigReqs.Map[sr.Pubkey] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: sr.SigningRoot}
		pubkeyToUUID[sr.Pubkey] = sr.UUID.String()
	}
	err := RequestValidatorSignaturesAsync(ctx, batchSigReqs, pubkeyToUUID)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("signType", sq.Type).Msg("RequestValidatorSignaturesAsync")
	}
}

func WaitForSignature(ctx context.Context, sr SignRequest) (Eth2SignResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
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
			log.Info().Interface("signRequest", sr).Interface("signResp", sigResp).Msg("found signature in cache")
			return Eth2SignResponse{Signature: sigResp}
		}
		time.Sleep(5 * time.Millisecond)
	}
}
