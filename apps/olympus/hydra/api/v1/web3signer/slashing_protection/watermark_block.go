package ethereum_slashing_protection_watermarking

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	dynamodb_web3signer "github.com/zeus-fyi/olympus/datastores/dynamodb/apps"
	dynamodb_web3signer_client "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/dynamodb_web3signer"
)

func WatermarkBlock(ctx context.Context, pubkey string, proposedSlot int) error {
	prevSlot, err := FetchLastSignedBlockSlot(ctx, pubkey)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to fetch last signed block slot")
		return err
	}
	if proposedSlot <= prevSlot {
		log.Ctx(ctx).Warn().Msgf("proposed slot %d less than or equal to a previous block slot %d", proposedSlot, prevSlot)
		return errors.New("proposed slot less than or equal to a previous block slot")
	}
	return nil
}

func FetchLastSignedBlockSlot(ctx context.Context, pubkey string) (prevSlot int, err error) {
	dynamoInstance := dynamodb_web3signer_client.Web3SignerDynamoDBClient
	key := dynamodb_web3signer.Web3SignerDynamoDBTableKeys{
		Pubkey:  pubkey,
		Network: Network,
	}
	bp, err := dynamoInstance.GetBlockProposal(ctx, key)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("key", key).Msg("failed to get last block proposal")
		return 0, err
	}
	return bp.Slot, nil
}
