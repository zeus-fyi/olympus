package hydra_eth2_web3signer

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
	ethereum_slashing_protection_watermarking "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/slashing_protection"
	bls_serverless_signing "github.com/zeus-fyi/zeus/pkg/aegis/aws/serverless_signing"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	"strconv"
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

func Watermarking(pubkey string, w *Web3SignerRequest) error {
	var sr SignRequest
	signingRoot := w.Body["signingRoot"]
	signType := w.Body["type"]

	sr.UUID = uuid.New()
	sr.Type = signType.(string)
	sr.Pubkey = pubkey
	sr.SigningRoot = signingRoot.(string)

	ctx := context.Background()
	switch signType {
	case ATTESTATION:
		log.Info().Interface("pubkey", pubkey).Msg("ATTESTATION")
		attestation := consensys_eth2_openapi.AttestationSigning{}
		b, err := json.Marshal(w.Body)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("ATTESTATION")
			return err
		}
		err = json.Unmarshal(b, &attestation)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("ATTESTATION")
			return err
		}
		// TODO watermark
		err = CanSignAttestation(ctx, pubkey, attestation.Attestation)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("ATTESTATION")
			return err
		}
	case AGGREGATION_SLOT:
		log.Info().Interface("pubkey", pubkey).Msg("AGGREGATION_SLOT")
	case AGGREGATE_AND_PROOF:
		log.Info().Interface("pubkey", pubkey).Msg("AGGREGATE_AND_PROOF")
	case BLOCK:
		log.Info().Interface("pubkey", pubkey).Msg("BLOCK")
		bs := consensys_eth2_openapi.BlockSigning{}
		b, err := json.Marshal(w.Body)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("BLOCK")
			return err
		}
		err = json.Unmarshal(b, &bs)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("BLOCK")
			return err
		}
		slot, err := strconv.Atoi(bs.Block.Slot)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Interface("beaconBody", bs.Block).Interface("slot", slot).Msg("BLOCK_V2")
			return err
		}
		err = ethereum_slashing_protection_watermarking.WatermarkBlock(ctx, pubkey, slot)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Interface("slot", slot).Msg("BLOCK")
			return err
		}
	case BLOCK_V2:
		log.Info().Interface("pubkey", pubkey).Msg("BLOCK_V2")
		bs := consensys_eth2_openapi.BeaconBlockSigning{}
		b, err := json.Marshal(w.Body)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("BLOCK_V2")
			return err
		}
		err = json.Unmarshal(b, &bs)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Msg("BLOCK_V2")
			return err
		}
		beaconBody, slot, err := DecodeBeaconBlockAndSlot(ctx, bs)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Interface("beaconBody", beaconBody).Interface("slot", slot).Msg("BLOCK_V2")
			return err
		}
		err = ethereum_slashing_protection_watermarking.WatermarkBlock(ctx, pubkey, slot)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("body", w.Body).Interface("slot", slot).Msg("BLOCK_V2")
			return err
		}
	case RANDAO_REVEAL:
		log.Info().Interface("pubkey", pubkey).Msg("RANDAO_REVEAL")
	case SYNC_COMMITTEE_MESSAGE:
		log.Info().Interface("pubkey", pubkey).Msg("SYNC_COMMITTEE_MESSAGE")
	case SYNC_COMMITTEE_SELECTION_PROOF:
		log.Info().Interface("pubkey", pubkey).Msg("SYNC_COMMITTEE_SELECTION_PROOF")
	case SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF:
		log.Info().Interface("pubkey", pubkey).Msg("SYNC_COMMITTEE_CONTRIBUTION_AND_PROOF")
	case VALIDATOR_REGISTRATION:
		log.Info().Interface("pubkey", pubkey).Msg("VALIDATOR_REGISTRATION")
	default:
	}

	// ideally can aggregate requests, and send in batch
	// TODO lookup secret name
	sigReqs := bls_serverless_signing.SignatureRequests{
		SecretName:        "",
		SignatureRequests: aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)},
	}

	// TODO add to stack -> broadcast to all
	sigReqs.SignatureRequests.Map[sr.Pubkey] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: sr.SigningRoot}
	return nil
}
