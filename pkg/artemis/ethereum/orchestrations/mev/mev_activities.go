package artemis_mev_transcations

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	mempool_txs "github.com/zeus-fyi/olympus/datastores/dynamodb/mempool"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

const hardhatSvc = "http://hardhat.hardhat.svc.cluster.local:8545"

func (d *ArtemisMevActivities) HistoricalSimulateAndValidateTx(ctx context.Context, trade artemis_autogen_bases.EthMempoolMevTx) error {
	wc := web3_client.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnetQuiknode.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnet.Account)
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
	wc := web3_client.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnetQuiknode.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnet.Account)
	txs, terr := wc.GetBlockTxs(ctx)
	if terr != nil {
		log.Err(terr).Str("network", d.Network).Msg("GetMempoolTxs failed")
		return terr
	}
	for _, tx := range txs {
		c.Set(tx.Hash().String(), tx, cache.DefaultExpiration)
		txBlackList := mempool_txs.TxBlacklistDynamoDB{
			TxBlacklistDynamoDBTableKeys: mempool_txs.TxBlacklistDynamoDBTableKeys{
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

func (d *ArtemisMevActivities) RemoveProcessedTx(ctx context.Context, tx mempool_txs.MempoolTxsDynamoDB) error {
	err := artemis_orchestration_auth.MevDynamoDBClient.RemoveMempoolTx(ctx, tx)
	if err != nil {
		log.Err(err).Str("network", d.Network).Msg("RemoveMempoolTx failed")
		return err
	}
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

func (d *ArtemisMevActivities) ConvertMempoolTxs(ctx context.Context, mempoolTxs []mempool_txs.MempoolTxsDynamoDB) (map[string]map[string]*types.Transaction, error) {
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
