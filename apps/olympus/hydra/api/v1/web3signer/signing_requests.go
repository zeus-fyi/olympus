package hydra_eth2_web3signer

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
	ethereum_slashing_protection_watermarking "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/slashing_protection"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

const (
	ATTESTATION                           = "ATTESTATION"
	AGGREGATION_SLOT                      = "AGGREGATION_SLOT"
	AGGREGATE_AND_PROOF                   = "AGGREGATE_AND_PROOF"
	BLOCK                                 = "BLOCK"
	BLOCK_V2                              = "BLOCK_V2"
	RANDAO_REVEAL                         = "RANDAO_REVEAL"
	SYNC_COMMITTEE_MESSAGE                = "SYNC_COMMITTEE_MESSAGE"
	SYNC_COMMITTEE_SELECTION_PROOF        = "SYNC_COMMITTEE_SELECTION_PROOF"
	SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF = "SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF"
	VALIDATOR_REGISTRATION                = "VALIDATOR_REGISTRATION"
)

type SignRequest struct {
	UUID uuid.UUID `json:"uuid"`

	Pubkey      string `json:"pubkey"`
	Type        string `json:"type"`
	SigningRoot string `json:"signingRoot"`
}

func Watermarking(ctx context.Context, pubkey string, w *Web3SignerRequest) (SignRequest, error) {
	var sr SignRequest
	signingRoot := w.Body["signingRoot"]
	signType := w.Body["type"]

	sr.UUID = uuid.New()
	sr.Type = signType.(string)
	sr.Pubkey = pubkey
	sr.SigningRoot = strings_filter.Trim0xPrefix(signingRoot.(string))
	log.Info().Interface("signType", signType).Interface("pubkey", sr.Pubkey).Msg("signType")

	switch signType {
	case ATTESTATION:
		log.Info().Interface("pubkey", pubkey).Msg("ATTESTATION")
		attestation := consensys_eth2_openapi.AttestationSigning{}
		b, err := json.Marshal(w.Body)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("ATTESTATION")
			return SignRequest{}, err
		}
		err = json.Unmarshal(b, &attestation)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("ATTESTATION")
			return SignRequest{}, err
		}
		err = CanSignAttestation(ctx, pubkey, attestation.Attestation)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("ATTESTATION")
			return SignRequest{}, err
		}
		AttestationSigningRequestPriorityQueue.PriorityQueue.Enqueue(sr)
	case AGGREGATION_SLOT:
		log.Info().Interface("pubkey", pubkey).Msg("AGGREGATION_SLOT")
		AggregationSlotSigningRequestPriorityQueue.PriorityQueue.Enqueue(sr)
	case AGGREGATE_AND_PROOF:
		log.Info().Interface("pubkey", pubkey).Msg("AGGREGATE_AND_PROOF")
		AggregationAndProofSigningRequestPriorityQueue.PriorityQueue.Enqueue(sr)
	case BLOCK:
		log.Info().Interface("pubkey", pubkey).Msg("BLOCK")
		bs := consensys_eth2_openapi.BlockSigning{}
		b, err := json.Marshal(w.Body)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("BLOCK")
			return SignRequest{}, err
		}
		err = json.Unmarshal(b, &bs)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("BLOCK")
			return SignRequest{}, err
		}
		slot, err := strconv.Atoi(bs.Block.Slot)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Interface("beaconBody", bs.Block).Interface("slot", slot).Msg("BLOCK_V2")
			return SignRequest{}, err
		}
		err = ethereum_slashing_protection_watermarking.WatermarkBlock(ctx, pubkey, slot)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Interface("slot", slot).Msg("BLOCK")
			return SignRequest{}, err
		}
		BlockSigningRequestPriorityQueue.PriorityQueue.Enqueue(sr)
	case BLOCK_V2:
		log.Info().Interface("pubkey", pubkey).Interface("body", w.Body).Msg("BLOCK_V2")
		bs := consensys_eth2_openapi.BeaconBlockSigning{}
		b, err := json.Marshal(w.Body)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("BLOCK_V2_MARSHALL_ERROR")
			return SignRequest{}, err
		}
		err = json.Unmarshal(b, &bs)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("BLOCK_V2_UNMARSHALL_ERROR")
			return SignRequest{}, err
		}
		beaconBody, slot, err := DecodeBeaconBlockAndSlot(ctx, bs.BeaconBlock)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Interface("beaconBody", beaconBody).Interface("SLOT", slot).Msg("BLOCK_V2_DECODE_ERROR")
			return SignRequest{}, err
		}
		err = ethereum_slashing_protection_watermarking.WatermarkBlock(ctx, pubkey, slot)
		if err != nil {
			log.Ctx(ctx).Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Interface("SLOT", slot).Msg("BLOCK_V2_WatermarkBlock_ERROR")
			return SignRequest{}, err
		}
		log.Info().Interface("pubkey", pubkey).Interface("body", w.Body).Interface("slot", slot).Msg("BLOCK_V2_WATERMARKED")
		BlockSigningRequestPriorityQueue.PriorityQueue.Enqueue(sr)
	case RANDAO_REVEAL:
		log.Info().Interface("pubkey", pubkey).Msg("RANDAO_REVEAL")
		RandaoRevealSigningRequestPriorityQueue.PriorityQueue.Enqueue(sr)
	case SYNC_COMMITTEE_MESSAGE:
		log.Info().Interface("pubkey", pubkey).Msg("SYNC_COMMITTEE_MESSAGE")
		SyncCommitteeMessageSigningRequestPriorityQueue.PriorityQueue.Enqueue(sr)
	case SYNC_COMMITTEE_SELECTION_PROOF:
		log.Info().Interface("pubkey", pubkey).Msg("SYNC_COMMITTEE_SELECTION_PROOF")
		SyncCommitteeSelectionProofSigningRequestPriorityQueue.PriorityQueue.Enqueue(sr)
	case SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF:
		log.Info().Interface("pubkey", pubkey).Msg("SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF")
		SyncCommitteeContributionAndProofPriorityQueue.PriorityQueue.Enqueue(sr)
	case VALIDATOR_REGISTRATION:
		log.Info().Interface("pubkey", pubkey).Msg("VALIDATOR_REGISTRATION")
		ValidatorRegistrationPriorityQueue.PriorityQueue.Enqueue(sr)
	default:
		log.Warn().Interface("pubkey", pubkey).Interface("type", w.Body["type"]).Interface("body", w.Body).Msg("UNKNOWN")
	}

	return sr, nil
}
