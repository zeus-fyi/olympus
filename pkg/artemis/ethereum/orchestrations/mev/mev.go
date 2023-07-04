package artemis_mev_transcations

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/apollo/ethereum/client_apis/beacon_api"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

const (
	irisSvc     = "https://iris.zeus.fyi/v1/internal/"
	irisBetaSvc = "https://iris.zeus.fyi/v1beta/internal/"
	hardhatSvc  = "https://hardhat.zeus.fyi/"
)

var (
	AuthHeader     string
	HardHatAccount *accounts.Account
)

func InitUniswap(ctx context.Context, authHeader string) {
	AuthHeader = authHeader
	newAccount, err := accounts.ParsePrivateKey("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		panic(err)
	}
	HardHatAccount = newAccount
	//go ProcessMempoolTxs(ctx)
}

func InitNewUniHardhat(ctx context.Context) *web3_client.UniswapClient {
	wc := web3_client.NewWeb3Client(irisBetaSvc, HardHatAccount)
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
	go beacon_api.TriggerWorkflowOnNewBlockHeaderEvent(ctx, artemis_network_cfgs.ArtemisQuicknodeStreamWebsocket, timestampChan)

	for {
		select {
		case t := <-timestampChan:
			log.Info().Msg(fmt.Sprintf("Received new timestamp: %s", t))
			log.Info().Msg("ExecuteArtemisMevWorkflow: ExecuteArtemisBlacklistTxWorkflow")
			err := ArtemisMevWorkerMainnet.ExecuteArtemisBlacklistTxWorkflow(ctx)
			log.Info().Msg("ExecuteArtemisMevWorkflow")
			time.Sleep(t.Add(8 * time.Second).Sub(time.Now()))
			err = ArtemisMevWorkerMainnet.ExecuteArtemisMevWorkflow(ctx)
			if err != nil {
				log.Err(err).Msg("ExecuteArtemisMevWorkflow failed")
			}
		}
	}
}
