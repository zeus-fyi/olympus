package ethereum_slashing_protection_watermarking

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"
	dynamodb_web3signer "github.com/zeus-fyi/olympus/datastores/dynamodb/apps"
	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
	dynamodb_web3signer_client "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/dynamodb_web3signer"
)

func WatermarkAttestation(ctx context.Context, pubkey string, attData consensys_eth2_openapi.AttestationData) error {
	sourceEpoch, targetEpoch, err := ConvertAttSourceTargetsToInt(ctx, attData)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to convert attestation source & target epoch to int")
		return err
	}
	prevSourceEpoch, prevTargetEpoch, prevAttestationHash, err := FetchLastAttestation(ctx, pubkey)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Msg("failed to fetch last attestation")
		return err
	}
	if IsSurroundVote(ctx, pubkey, prevSourceEpoch, prevTargetEpoch, sourceEpoch, targetEpoch) {
		err = errors.New("surround vote")
		log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Msg("surround vote")
		return err
	}
	newAttestationDataHash, err := HashAttestationData(attData)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Msg("failed to hash attestation data")
		return err
	}
	// double vote (data_1 != data_2 and data_1.target.epoch == data_2.target.epoch)
	if prevTargetEpoch == targetEpoch && prevAttestationHash != newAttestationDataHash {
		err = errors.New("double vote")
		log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Interface("prevAttestationHash", prevAttestationHash).Interface("newAttestationDataHash", newAttestationDataHash).Msg("double vote")
		return err
	}
	if sourceEpoch < prevSourceEpoch {
		err = errors.New("source epoch less than previously recorded epoch")
		log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey)
		return err
	}
	if sourceEpoch > targetEpoch {
		err = errors.New("source epoch can't be greater than target epoch")
		log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey)
		return err
	}
	if targetEpoch < prevTargetEpoch {
		err = errors.New("new target epoch can't be less than previously recorded target epoch")
		log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey)
		return err
	}
	return RecordAttestation(ctx, pubkey, newAttestationDataHash, sourceEpoch, targetEpoch)
}

func HashAttestationData(attData consensys_eth2_openapi.AttestationData) (string, error) {
	// Encode the attestation data as JSON
	bytes, err := json.Marshal(attData)
	if err != nil {
		return "", err
	}
	// Compute the SHA256 hash of the JSON bytes
	hash := sha256.Sum256(bytes)
	hexStr := hex.EncodeToString(hash[:])
	return hexStr, nil
}
func RecordAttestation(ctx context.Context, pubkey, attData string, sourceEpoch, targetEpoch int) error {
	dynamoInstance := dynamodb_web3signer_client.Web3SignerDynamoDBClient
	key := dynamodb_web3signer.Web3SignerDynamoDBTableKeys{
		Pubkey:  pubkey,
		Network: Network,
	}
	dynAtt := dynamodb_web3signer.AttestationsDynamoDB{
		Web3SignerDynamoDBTableKeys: key,
		SourceEpoch:                 sourceEpoch,
		TargetEpoch:                 targetEpoch,
		AttestationData:             attData,
	}
	err := dynamoInstance.PutAttestation(ctx, dynAtt)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("key", key).Msg("failed to get attestation")
		return err
	}
	return nil
}

func FetchLastAttestation(ctx context.Context, pubkey string) (sourceEpoch, targetEpoch int, attestationHash string, err error) {
	dynamoInstance := dynamodb_web3signer_client.Web3SignerDynamoDBClient
	key := dynamodb_web3signer.Web3SignerDynamoDBTableKeys{
		Pubkey:  pubkey,
		Network: Network,
	}
	att, err := dynamoInstance.GetAttestation(ctx, key)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("key", key).Msg("failed to get attestation")
		return 0, 0, "", err
	}
	return att.SourceEpoch, att.TargetEpoch, att.AttestationData, nil
}

func IsSurroundVote(ctx context.Context, pubkey string, data1SourceEpoch, data1TargetEpoch, data2SourceEpoch, data2TargetEpoch int) bool {
	if data1SourceEpoch < data2SourceEpoch && data2TargetEpoch < data1TargetEpoch {
		log.Ctx(ctx).Warn().Interface("data1SourceEpoch", data1SourceEpoch).Interface("data1TargetEpoch", data1TargetEpoch).Interface("data2SourceEpoch", data2SourceEpoch).Interface("data2TargetEpoch", data2TargetEpoch).Msgf("detected surround vote for %s", pubkey)
		return true
	}
	return false
}
