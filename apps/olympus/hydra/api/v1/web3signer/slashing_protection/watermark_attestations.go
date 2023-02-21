package ethereum_slashing_protection_watermarking

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	dynamodb_web3signer "github.com/zeus-fyi/olympus/datastores/dynamodb/apps"
	dynamodb_web3signer_client "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/dynamodb_web3signer"
)

func WatermarkAttestation(ctx context.Context, pubkey string, sourceEpoch, targetEpoch int) error {
	prevSourceEpoch, prevTargetEpoch, err := FetchLastAttestation(ctx, pubkey)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Msg("failed to fetch last attestation")
		return err
	}
	if IsSurroundVote(ctx, pubkey, prevSourceEpoch, prevTargetEpoch, sourceEpoch, targetEpoch) {
		log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Msg("surround vote")
		err = errors.New("surround vote")
		return err
	}
	if IsSourceEpochLessThanOrEqualToAnyPreviousAttestations(ctx, pubkey, sourceEpoch, prevSourceEpoch) {
		return errors.New("sourceEpoch less than or equal to a previous attestation sourceEpoch")
	}
	if IsTargetEpochLessThanOrEqualToAnyPreviousAttestations(ctx, pubkey, targetEpoch, prevTargetEpoch) {
		return errors.New("targetEpoch less than or equal to a previous attestation targetEpoch")
	}
	return RecordAttestation(ctx, pubkey, sourceEpoch, targetEpoch)
}

func RecordAttestation(ctx context.Context, pubkey string, sourceEpoch, targetEpoch int) error {
	dynamoInstance := dynamodb_web3signer_client.Web3SignerDynamoDBClient
	key := dynamodb_web3signer.Web3SignerDynamoDBTableKeys{
		Pubkey:  pubkey,
		Network: Network,
	}
	dynAtt := dynamodb_web3signer.AttestationsDynamoDB{
		Web3SignerDynamoDBTableKeys: key,
		SourceEpoch:                 sourceEpoch,
		TargetEpoch:                 targetEpoch,
	}
	err := dynamoInstance.PutAttestation(ctx, dynAtt)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("key", key).Msg("failed to get attestation")
		return err
	}
	return nil
}

func FetchLastAttestation(ctx context.Context, pubkey string) (sourceEpoch, targetEpoch int, err error) {
	dynamoInstance := dynamodb_web3signer_client.Web3SignerDynamoDBClient

	key := dynamodb_web3signer.Web3SignerDynamoDBTableKeys{
		Pubkey:  pubkey,
		Network: Network,
	}
	att, err := dynamoInstance.GetAttestation(ctx, key)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("key", key).Msg("failed to get attestation")
		return 0, 0, err
	}
	return att.SourceEpoch, att.TargetEpoch, nil
}

// TODO far future signing protection

// IsSourceEpochGreaterThanTargetEpoch if this is true then reject the signing request
func IsSourceEpochGreaterThanTargetEpoch(ctx context.Context, pubkey string, sourceEpoch, targetEpoch int) bool {
	if sourceEpoch >= targetEpoch {
		log.Ctx(ctx).Warn().Msgf("detected sourceEpoch %d greater than or equal targetEpoch %d for %s", sourceEpoch, targetEpoch, pubkey)
		return true
	}
	return false
}

// IsSourceEpochLessThanOrEqualToAnyPreviousAttestations if this is true then reject the signing request. hasSourceOlderThanWatermark
func IsSourceEpochLessThanOrEqualToAnyPreviousAttestations(ctx context.Context, pubkey string, newSourceEpoch, prevSourceEpoch int) bool {
	if newSourceEpoch <= prevSourceEpoch {
		log.Ctx(ctx).Warn().Msgf("detected newSourceEpoch %d less than or equal to maxRecordedSourceEpoch %d for %s", newSourceEpoch, prevSourceEpoch, pubkey)
		return true
	}
	return false
}

func IsTargetEpochLessThanOrEqualToAnyPreviousAttestations(ctx context.Context, pubkey string, newTargetEpoch, prevTargetEpoch int) bool {
	if newTargetEpoch <= prevTargetEpoch {
		log.Ctx(ctx).Warn().Msgf("detected newSourceEpoch %d less than or equal to maxRecordedSourceEpoch %d for %s", newTargetEpoch, prevTargetEpoch, pubkey)
		return true
	}
	return false
}

func IsSurroundVote(ctx context.Context, pubkey string, data1SourceEpoch, data1TargetEpoch, data2SourceEpoch, data2TargetEpoch int) bool {
	if data1SourceEpoch < data2SourceEpoch && data2TargetEpoch < data1TargetEpoch {
		log.Ctx(ctx).Warn().Interface("data1SourceEpoch", data1SourceEpoch).Interface("data1TargetEpoch", data1TargetEpoch).Interface("data2SourceEpoch", data2SourceEpoch).Interface("data2TargetEpoch", data2TargetEpoch).Msgf("detected surround vote for %s", pubkey)
		return true
	}
	return false
}
