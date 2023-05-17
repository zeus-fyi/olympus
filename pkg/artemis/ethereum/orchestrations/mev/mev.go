package artemis_mev_tx_fetcher

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
	go GetMempoolTxs(ctx)
}

func GetMempoolTxs(ctx context.Context) {
	for {
		go func() {
			txMap, err := Uniswap.Web3Client.GetFilteredPendingMempoolTxs(ctx, Uniswap.MevSmartContractTxMap)
			if err != nil {
				log.Error().Err(err).Msg("failed to get mempool txs")
			}
			Uniswap.MevSmartContractTxMap = txMap
			Uniswap.ProcessTxs(ctx)
		}()
		time.Sleep(100 * time.Millisecond)
	}
}
