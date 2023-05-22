package artemis_mev_transcations

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

var Uniswap web3_client.UniswapV2Client

func InitUniswap(ctx context.Context, authHeader string) {
	wc := web3_client.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnet.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnet.Account)

	m := map[string]string{
		"Authorization": "Bearer " + authHeader,
	}
	wc.Headers = m
	Uniswap = web3_client.InitUniswapV2Client(ctx, wc)
	Uniswap.PrintOn = true
	Uniswap.PrintLocal = false
	go ProcessMempoolTxs(ctx)
}

func ProcessMempoolTxs(ctx context.Context) {
	for {
		err := ArtemisMevWorkerMainnet.ExecuteArtemisMevWorkflow(ctx)
		if err != nil {
			log.Err(err).Msg("ExecuteArtemisMevWorkflow failed")
		}
		time.Sleep(100000 * time.Millisecond)
	}
}
