package artemis_mev_transcations

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/apollo/ethereum/client_apis/beacon_api"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

const (
	irisSvc                = "https://iris.zeus.fyi/v2/internal/router"
	irisBetaSvc            = "https://iris.zeus.fyi/v2/internal/router"
	hardhatSvc             = "https://hardhat.zeus.fyi/"
	irisSvcBeaconsInternal = "http://iris.iris.svc.cluster.local/v2/internal/router/"
	irisBetaSvcInternal    = "http://iris.iris.svc.cluster.local/v2/internal/router/"
)

var (
	AuthHeader string
)

func InitArtemisUniswap(ctx context.Context, authHeader string) {
	AuthHeader = authHeader
	timestampChan := make(chan time.Time)
	go ProcessMempoolTxs(context.Background(), timestampChan)
	go artemis_trading_cache.SetActiveTradingBlockCache(context.Background(), timestampChan)
	go beacon_api.TriggerWorkflowOnNewBlockHeaderEvent(context.Background(), artemis_network_cfgs.ArtemisQuicknodeStreamWebsocket, timestampChan)
}

func InitTycheUniswap(ctx context.Context, authHeader string) {
}
func InitNewUniHardhat(ctx context.Context, sessionID string) *web3_client.UniswapClient {
	acc, err := accounts.CreateAccount()
	if err != nil {
		panic(err)
	}
	wc := web3_client.NewWeb3Client(irisSvcBeaconsInternal, acc)
	wc.Network = hestia_req_types.Mainnet
	wc.AddDefaultEthereumMainnetTableHeader()
	wc.AddMaxBlockHeightProcedureEthJsonRpcHeader()
	wc.AddBearerToken(AuthHeader)
	wc.AddSessionLockHeader(sessionID)
	uni := web3_client.InitUniswapClient(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	uni.Web3Client.IsAnvilNode = true
	uni.Web3Client.DurableExecution = true
	return &uni
}

func InitNewUniswapQuikNode(ctx context.Context) *web3_client.UniswapClient {
	wc := web3_client.NewWeb3Client(irisSvcBeaconsInternal, artemis_network_cfgs.ArtemisEthereumMainnet.Account)
	wc.Network = hestia_req_types.Mainnet
	uni := web3_client.InitUniswapClient(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	return &uni
}

func ProcessMempoolTxs(ctx context.Context, timestampChan chan time.Time) {

	for {
		select {
		case t := <-timestampChan:
			// todo: margin needed for now, but should remove via refactor
			time.Sleep(25 * time.Millisecond)
			wc := web3_actions.NewWeb3ActionsClient(irisSvcBeaconsInternal)
			wc.AddDefaultEthereumMainnetTableHeader()

			wc.Network = hestia_req_types.Mainnet
			wc.Dial()
			bn, berr := artemis_trading_cache.GetLatestBlockFromCacheOrProvidedSource(context.Background(), wc)
			if berr != nil {
				log.Err(berr).Msg("failed to get block number")
				wc.Close()
			}
			wc.Close()
			err := ArtemisMevWorkerMainnet2.ExecuteArtemisGetLookaheadPricesWorkflow(context.Background(), bn)
			if err != nil {
				log.Err(err).Msg("ExecuteArtemisMevWorkflow failed")
			}
			log.Info().Msg(fmt.Sprintf("Received new timestamp: %s", t))
			log.Info().Msg("ExecuteArtemisMevWorkflow: ExecuteArtemisBlacklistTxWorkflow")
			err = ArtemisMevWorkerMainnet.ExecuteArtemisBlacklistTxWorkflow(context.Background())
			if err != nil {
				log.Err(err).Msg("ExecuteArtemisBlacklistTxWorkflow failed")
			}
			log.Info().Msg("ExecuteArtemisMevWorkflow")
			time.Sleep(t.Add(8 * time.Second).Sub(time.Now()))
			err = ArtemisMevWorkerMainnetHistoricalTxs.ExecuteArtemisMevWorkflow(context.Background(), int(bn))
			if err != nil {
				log.Err(err).Msg("ExecuteArtemisMevWorkflow failed")
			}
		}
	}
}
