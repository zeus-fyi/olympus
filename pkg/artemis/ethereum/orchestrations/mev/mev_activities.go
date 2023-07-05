package artemis_mev_transcations

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	dynamodb_mev "github.com/zeus-fyi/olympus/datastores/dynamodb/mev"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (d *ArtemisMevActivities) HistoricalSimulateAndValidateTx(ctx context.Context, trade artemis_autogen_bases.EthMempoolMevTx) error {
	wc := web3_client.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalPrimary.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnet.Account)
	uni := InitNewUniHardhat(ctx)
	err := uni.RunHistoricalTradeAnalysis(ctx, trade.TxFlowPrediction, wc)
	if err != nil {
		log.Err(err).Msg("RunHistoricalTradeAnalysis failed")
		return err
	}
	uni.PrintResults()
	return nil
}

func (d *ArtemisMevActivities) SimulateAndValidateBundle(ctx context.Context) error {
	return nil
}

func (d *ArtemisMevActivities) SubmitFlashbotsBundle(ctx context.Context) error {
	return nil
}

var c = cache.New(5*time.Hour, 10*time.Hour)

func (d *ArtemisMevActivities) BlacklistMinedTxs(ctx context.Context) error {
	wc := web3_client.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnet.Account)
	txs, terr := wc.GetBlockTxs(ctx)
	if terr != nil {
		log.Err(terr).Str("network", d.Network).Msg("GetDynamoDBMempoolTxs failed")
		return terr
	}
	for _, tx := range txs {
		c.Set(tx.Hash().String(), tx, cache.DefaultExpiration)
		txBlackList := dynamodb_mev.TxBlacklistDynamoDB{
			TxBlacklistDynamoDBTableKeys: dynamodb_mev.TxBlacklistDynamoDBTableKeys{
				TxHash: tx.Hash().String(),
			},
		}
		err := artemis_orchestration_auth.MevDynamoDBClient.PutTxBlacklist(ctx, txBlackList)
		if err != nil {
			log.Err(err).Str("network", d.Network).Msg("PutTxBlacklist failed")
			return err
		}
	}
	return nil
}

func (d *ArtemisMevActivities) BlacklistProcessedTxs(ctx context.Context, txSlice artemis_autogen_bases.EthMempoolMevTxSlice) error {
	for _, tx := range txSlice {
		c.Set(tx.TxHash, tx, cache.DefaultExpiration)
		txBlackList := dynamodb_mev.TxBlacklistDynamoDB{
			TxBlacklistDynamoDBTableKeys: dynamodb_mev.TxBlacklistDynamoDBTableKeys{
				TxHash: tx.TxHash,
			},
		}
		err := artemis_orchestration_auth.MevDynamoDBClient.PutTxBlacklist(ctx, txBlackList)
		if err != nil {
			log.Err(err).Str("network", d.Network).Msg("BlacklistProcessedTxs failed")
			return err
		}
	}
	return nil
}

func (d *ArtemisMevActivities) RemoveProcessedTx(ctx context.Context, tx dynamodb_mev.MempoolTxsDynamoDB) error {
	err := artemis_orchestration_auth.MevDynamoDBClient.RemoveMempoolTx(ctx, tx)
	if err != nil {
		log.Err(err).Str("network", d.Network).Msg("RemoveMempoolTx failed")
		return err
	}
	return nil
}

func (d *ArtemisMevActivities) GetDynamoDBMempoolTxs(ctx context.Context) ([]dynamodb_mev.MempoolTxsDynamoDB, error) {
	txs, terr := artemis_orchestration_auth.MevDynamoDBClient.GetMempoolTxs(ctx, d.Network)
	if terr != nil {
		log.Err(terr).Str("network", d.Network).Msg("GetDynamoDBMempoolTxs failed")
		return nil, terr
	}
	return txs, nil
}

func (d *ArtemisMevActivities) GetPostgresMempoolTxs(ctx context.Context, bn int) (artemis_autogen_bases.EthMempoolMevTxSlice, error) {
	txs, terr := artemis_mev_models.SelectMempoolTxAtBlockNumber(ctx, 1, bn)
	if terr != nil {
		log.Err(terr).Str("network", d.Network).Msg("GetPostgresMempoolTxs failed")
		return nil, terr
	}
	return txs, nil
}

func (d *ArtemisMevActivities) ConvertMempoolTxs(ctx context.Context, mempoolTxs []dynamodb_mev.MempoolTxsDynamoDB) (map[string]map[string]*types.Transaction, error) {
	txMap := make(map[string]map[string]*types.Transaction)
	for _, tx := range mempoolTxs {
		if txMap[tx.Pubkey] == nil {
			txMap[tx.Pubkey] = make(map[string]*types.Transaction)
		}
		txRpc := &types.Transaction{}
		b, berr := json.Marshal(tx.Tx)
		if berr != nil {
			log.Err(berr).Msg("ConvertMempoolTxs: error marshalling tx")
			return nil, berr
		}
		var txRpcMapStr string
		berr = json.Unmarshal(b, &txRpcMapStr)
		if berr != nil {
			log.Err(berr).Msg("ConvertMempoolTxs: error marshalling tx")
			return nil, berr
		}
		berr = json.Unmarshal([]byte(txRpcMapStr), &txRpc)
		if berr != nil {
			log.Err(berr).Msg("ConvertMempoolTxs: error marshalling tx")
			return nil, berr
		}
		_, found := c.Get(txRpc.Hash().String())
		if found {
			log.Info().Str("txHash", txRpc.Hash().String()).Msg("tx already mined")
			continue
		}
		tmp := txMap[tx.Pubkey]
		tmp[fmt.Sprintf("%d", tx.TxOrder)] = txRpc
		txMap[tx.Pubkey] = tmp
	}
	return txMap, nil
}

func (d *ArtemisMevActivities) ProcessMempoolTxs(ctx context.Context, mempoolTxs map[string]map[string]*types.Transaction) ([]artemis_autogen_bases.EthMempoolMevTx, error) {
	uni := InitNewUniswapQuiknode(ctx)
	err := uni.ProcessMempoolTxs(ctx, mempoolTxs)
	if err != nil {
		log.Err(err).Msg("ProcessMempoolTxs failed")
		return nil, err
	}
	uni.ProcessTxs(ctx)
	return uni.Trades, nil
}
