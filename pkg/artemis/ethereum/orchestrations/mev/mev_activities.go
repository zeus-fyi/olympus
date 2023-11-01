package artemis_mev_transcations

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	dynamodb_mev "github.com/zeus-fyi/olympus/datastores/dynamodb/mev"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trade_debugger "github.com/zeus-fyi/olympus/pkg/artemis/trading/trade_debugger"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

func (d *ArtemisMevActivities) HistoricalSimulateAndValidateTx(ctx context.Context, trade artemis_autogen_bases.EthMempoolMevTx) (string, error) {
	sessionID := uuid.New().String()
	uni := InitNewUniHardhat(ctx, sessionID)
	at := artemis_realtime_trading.NewActiveTradingDebugger(uni)
	td := artemis_trade_debugger.NewTradeDebuggerWorkflowAnalysis(at, uni.Web3Client)
	err := td.Replay(context.Background(), trade.TxHash, true)
	if err != nil {
		log.Err(err).Str("sessionID", uni.Web3Client.GetSessionLockHeader()).Str("network", d.Network).Msg("Replay failed")
		return sessionID, err
	}
	return sessionID, err
}

func (d *ArtemisMevActivities) EndServerlessSession(ctx context.Context, sessionID string) error {
	irisClient := resty_base.GetBaseRestyClient("https://iris.zeus.fyi/v1/serverless", artemis_orchestration_auth.Bearer)
	resp, err := irisClient.R().
		Delete(fmt.Sprintf("/v1/serverless/%s", sessionID))
	if err != nil {
		log.Err(err).Str("network", d.Network).Msg("EndServerlessSession failed")
		return err
	}
	if resp.StatusCode() >= 400 {
		err = fmt.Errorf("EndServerlessSession: status code: %d", resp.StatusCode())
		log.Err(err).Str("network", d.Network).Msg("EndServerlessSession failed")
		return err
	}
	return nil
}

/*
	irisClient := resty_base.GetBaseRestyClient("http://localhost:9010", s.Tc.ProductionLocalTemporalBearerToken)

	sessionID := "sessionID"
	resp, err := irisClient.R().
		Delete(fmt.Sprintf("/v1/serverless/%s", sessionID))

*/
// end block number

func (d *ArtemisMevActivities) SimulateAndValidateBundle(ctx context.Context) error {
	return nil
}

func (d *ArtemisMevActivities) SubmitFlashbotsBundle(ctx context.Context) error {
	return nil
}

var c = cache.New(5*time.Hour, 10*time.Hour)

func (d *ArtemisMevActivities) BlacklistMinedTxs(ctx context.Context) error {
	wc := web3_client.NewWeb3Client(irisSvcBeaconsInternal, artemis_network_cfgs.ArtemisEthereumMainnet.Account)
	wc.AddBearerToken(artemis_orchestration_auth.Bearer)
	txs, terr := wc.GetBlockTxs(ctx)
	if terr != nil {
		log.Err(terr).Str("network", d.Network).Msg("BlacklistMinedTxs GetBlockTxs failed")
		return terr
	}
	for _, tx := range txs {
		c.Set(tx.Hash().String(), tx, cache.DefaultExpiration)
		err := artemis_trading_cache.WriteRedis.AddTxHashCache(context.Background(), tx.Hash().String(), time.Hour*24)
		if err != nil {
			log.Err(err).Str("network", d.Network).Msg("BlacklistMinedTxs: AddTxHashCache failed")
			return err
		}
	}
	return nil
}

func (d *ArtemisMevActivities) GetLookaheadPrices(ctx context.Context, bn uint64) error {
	wc := web3_actions.NewWeb3ActionsClient(irisSvcBeaconsInternal)
	wc.AddBearerToken(artemis_orchestration_auth.Bearer)
	wc.AddDefaultEthereumMainnetTableHeader()
	err := artemis_uniswap_pricing.FetchV2PairsToMulticall(ctx, wc, bn)
	if err != nil {
		return err
	}
	return nil
}

func (d *ArtemisMevActivities) BlacklistProcessedTxs(ctx context.Context, txSlice artemis_autogen_bases.EthMempoolMevTxSlice) error {
	for _, tx := range txSlice {
		c.Set(tx.TxHash, tx, cache.DefaultExpiration)
		err := artemis_trading_cache.WriteRedis.AddTxHashCache(context.Background(), tx.TxHash, time.Hour*24)
		if err != nil {
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
	uni := InitNewUniswapQuikNode(ctx)
	err := uni.ProcessMempoolTxs(ctx, mempoolTxs)
	if err != nil {
		log.Err(err).Msg("ProcessMempoolTxs failed")
		return nil, err
	}
	uni.ProcessTxs(ctx, nil)
	return uni.Trades, nil
}
