package artemis_mev_transcations

import (
	"context"

	"github.com/rs/zerolog/log"
	mempool_txs "github.com/zeus-fyi/olympus/datastores/dynamodb/mempool"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
)

func (d *ArtemisMevActivities) SimulateAndValidateBundle(ctx context.Context) error {
	return nil
}

func (d *ArtemisMevActivities) SubmitFlashbotsBundle(ctx context.Context) error {
	return nil
}

func (d *ArtemisMevActivities) GetMempoolTxs(ctx context.Context) ([]mempool_txs.MempoolTxsDynamoDB, error) {
	txs, err := artemis_orchestration_auth.MevDynamoDBClient.GetMempoolTxs(ctx, d.Network)
	if err != nil {
		log.Err(err).Str("network", d.Network).Msg("GetMempoolTxs failed")
	}
	return txs, err
}

func (d *ArtemisMevActivities) DecodeMempoolTxs(ctx context.Context) error {
	return nil
}
