package artemis_mev_transcations

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

var (
	AuthHeader string
)

func InitUniswap(ctx context.Context, authHeader string) {
	AuthHeader = authHeader
	go ProcessMempoolTxs(ctx)
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

func ProcessMempoolTxs(ctx context.Context) {
	for {
		err := ArtemisMevWorkerMainnet.ExecuteArtemisMevWorkflow(ctx)
		if err != nil {
			log.Err(err).Msg("ExecuteArtemisMevWorkflow failed")
		}
		time.Sleep(200 * time.Millisecond)
	}
}
