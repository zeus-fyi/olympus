package artemis_mev_transcations

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (d *ArtemisMevActivities) SimulateAndValidateBundle(ctx context.Context) error {
	return nil
}

func (d *ArtemisMevActivities) SubmitFlashbotsBundle(ctx context.Context) error {
	return nil
}

func (d *ArtemisMevActivities) GetMempoolTxs(ctx context.Context) (map[string]map[string]*types.Transaction, error) {
	txs, terr := artemis_orchestration_auth.MevDynamoDBClient.GetMempoolTxs(ctx, d.Network)
	if terr != nil {
		log.Err(terr).Str("network", d.Network).Msg("GetMempoolTxs failed")
		return nil, terr
	}
	return txs, nil
}

func (d *ArtemisMevActivities) ProcessMempoolTxs(ctx context.Context, mempoolTxs map[string]map[string]*types.Transaction) ([]artemis_autogen_bases.EthMempoolMevTx, error) {
	uni := InitNewUniswapQuiknode(ctx)
	mevTxMap := uni.MevSmartContractTxMap
	processedMevTxMap, err := web3_client.ProcessMempoolTxs(ctx, mempoolTxs, mevTxMap)
	if err != nil {
		log.Err(err).Msg("ProcessMempoolTxs failed")
		return nil, err
	}
	uni.MevSmartContractTxMap = processedMevTxMap
	uni.ProcessTxs(ctx)
	return uni.Trades, nil
}
