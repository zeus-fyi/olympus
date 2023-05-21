package artemis_mev_transcations

import (
	"context"
	"encoding/json"

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
	txs, terr := artemis_orchestration_auth.MevDynamoDBClient.GetMempoolTxs(ctx, d.Network)
	if terr != nil {
		log.Err(terr).Str("network", d.Network).Msg("GetMempoolTxs failed")
		return nil, terr
	}
	return txs, nil
}

func (d *ArtemisMevActivities) DecodeMempoolTxs(ctx context.Context, txs []mempool_txs.MempoolTxsDynamoDB) error {
	// TODO filter and process w/uniswap
	for _, tx := range txs {
		b, err := json.Marshal(tx.Tx)
		if err != nil {
			return err
		}
		txIn := map[string]interface{}{}
		err = json.Unmarshal(b, &txIn)
		if err != nil {
			return err
		}
	}
	return nil
}
