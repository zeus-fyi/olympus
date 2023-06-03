package artemis_mev_transcations

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
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
	go ProcessMempoolTxs(ctx)
}

func InitNewUniHardhat(ctx context.Context) *web3_client.UniswapV2Client {
	wc := web3_client.NewWeb3Client(hardhatSvc, HardHatAccount)
	m := map[string]string{
		"Authorization": "Bearer " + AuthHeader,
	}
	wc.Headers = m
	uni := web3_client.InitUniswapV2Client(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	return &uni
}

func InitNewUniswap(ctx context.Context) *web3_client.UniswapV2Client {
	wc := web3_client.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnet.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnet.Account)

	m := map[string]string{
		"Authorization": "Bearer " + AuthHeader,
	}
	wc.Headers = m
	uni := web3_client.InitUniswapV2Client(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	return &uni
}

func InitNewUniswapQuiknode(ctx context.Context) *web3_client.UniswapV2Client {
	wc := web3_client.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnetQuiknode.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnet.Account)
	uni := web3_client.InitUniswapV2Client(ctx, wc)
	uni.PrintOn = true
	uni.PrintLocal = false
	return &uni
}

func ProcessMempoolTxs(ctx context.Context) {
	cr := chronos.Chronos{}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			secondsLeftInSlot := 12 - cr.GetSecsSinceLastMainnetSlot()
			if secondsLeftInSlot <= 4 {
				// when 4 seconds remaining execute this
				log.Info().Msg("ExecuteArtemisMevWorkflow")
				err := ArtemisMevWorkerMainnet.ExecuteArtemisMevWorkflow(ctx)
				if err != nil {
					log.Err(err).Msg("ExecuteArtemisMevWorkflow failed")
				}
				time.Sleep(4 * time.Second)
				err = ArtemisMevWorkerMainnet.ExecuteArtemisBlacklistTxWorkflow(ctx)
				if err != nil {
					log.Err(err).Msg("ExecuteArtemisBlacklistTxWorkflow failed")
				}
			}

			if secondsLeftInSlot <= 0 {
				// when 0 seconds remaining go to start of loop
				secondsLeftInSlot = 12 - cr.GetSecsSinceLastMainnetSlot()
				ticker.Reset(time.Duration(secondsLeftInSlot) * time.Second)
			}
		}
	}
}
