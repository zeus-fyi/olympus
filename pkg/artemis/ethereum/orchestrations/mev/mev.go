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
)

const (
	irisSvc     = "https://iris.zeus.fyi/v1/internal/"
	irisBetaSvc = "https://iris.zeus.fyi/v1beta/internal/"
	hardhatSvc  = "https://hardhat.zeus.fyi/"
)

var (
	AuthHeader string
)

func InitUniswap(ctx context.Context, authHeader string) {
	AuthHeader = authHeader
	go ProcessMempoolTxs(ctx)
}

func InitNewUniHardhat(ctx context.Context) *web3_client.UniswapClient {
	acc, err := accounts.CreateAccount()
	if err != nil {
		panic(err)
	}
	wc := web3_client.NewWeb3Client(irisBetaSvc, acc)
	m := map[string]string{
		"Authorization": "Bearer " + AuthHeader,
	}
	wc.Headers = m
	uni := web3_client.InitUniswapClient(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	uni.Web3Client.IsAnvilNode = true
	uni.Web3Client.DurableExecution = true
	return &uni
}

func InitNewUniswap(ctx context.Context) *web3_client.UniswapClient {
	wc := web3_client.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnet.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnet.Account)

	m := map[string]string{
		"Authorization": "Bearer " + AuthHeader,
	}
	wc.Headers = m
	uni := web3_client.InitUniswapClient(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	return &uni
}

func InitNewUniswapQuiknode(ctx context.Context) *web3_client.UniswapClient {
	wc := web3_client.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnetQuiknode.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnet.Account)
	uni := web3_client.InitUniswapClient(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	return &uni
}

func ProcessMempoolTxs(ctx context.Context) {
	timestampChan := make(chan time.Time)
	go artemis_trading_cache.SetActiveTradingBlockCache(ctx)
	go beacon_api.TriggerWorkflowOnNewBlockHeaderEvent(ctx, artemis_network_cfgs.ArtemisQuicknodeStreamWebsocket, timestampChan)

	for {
		select {
		case t := <-timestampChan:
			// todo: margin needed for now, but should remove via refactor
			time.Sleep(25 * time.Millisecond)
			wc := web3_actions.NewWeb3ActionsClient(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL)
			wc.Dial()
			bn, berr := artemis_trading_cache.GetLatestBlockFromCacheOrProvidedSource(ctx, wc)
			if berr != nil {
				log.Err(berr).Msg("failed to get block number")
				return
			}
			wc.Close()
			err := ArtemisActiveMevWorkerMainnet.ExecuteArtemisGetLookaheadPricesWorkflow(ctx, bn)
			if err != nil {
				log.Err(err).Msg("ExecuteArtemisMevWorkflow failed")
			}
			log.Info().Msg(fmt.Sprintf("Received new timestamp: %s", t))
			log.Info().Msg("ExecuteArtemisMevWorkflow: ExecuteArtemisBlacklistTxWorkflow")
			err = ArtemisActiveMevWorkerMainnet.ExecuteArtemisBlacklistTxWorkflow(ctx)
			if err != nil {
				log.Err(err).Msg("ExecuteArtemisBlacklistTxWorkflow failed")
			}
			log.Info().Msg("ExecuteArtemisMevWorkflow")
			time.Sleep(t.Add(8 * time.Second).Sub(time.Now()))
			err = ArtemisMevWorkerMainnetHistoricalTxs.ExecuteArtemisMevWorkflow(ctx, int(bn))
			if err != nil {
				log.Err(err).Msg("ExecuteArtemisMevWorkflow failed")
			}
		}
	}
}
