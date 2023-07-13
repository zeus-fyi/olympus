package artemis_trading_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/apollo/ethereum/client_apis/beacon_api"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
)

var (
	TokenMap map[string]artemis_autogen_bases.Erc20TokenInfo
	Cache    = cache.New(12*time.Second, 4*time.Second)
)

func InitTokenFilter(ctx context.Context) {
	_, tm, terr := artemis_mev_models.SelectERC20Tokens(ctx)
	if terr != nil {
		panic(terr)
	}
	TokenMap = tm
}

var wc = web3_actions.NewWeb3ActionsClient(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL)

func GetLatestBlock(ctx context.Context) (uint64, error) {
	val, ok := Cache.Get("block_number")
	if ok && val != nil {
		return val.(uint64), nil
	}
	wc.Dial()
	defer wc.Close()
	bn, berr := wc.C.BlockNumber(ctx)
	if berr != nil {
		log.Err(berr).Msg("failed to get block number")
		return 0, berr
	}
	Cache.Set("block_number", bn, 12*time.Second)
	return bn, nil
}

func SetActiveTradingBlockCache(ctx context.Context) {
	timestampChan := make(chan time.Time)
	go beacon_api.TriggerWorkflowOnNewBlockHeaderEvent(ctx, artemis_network_cfgs.ArtemisQuicknodeStreamWebsocket, timestampChan)

	for {
		select {
		case t := <-timestampChan:
			Cache.Delete("block_number")

			wc = web3_actions.NewWeb3ActionsClient(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL)
			wc.Dial()
			bn, berr := wc.C.BlockNumber(ctx)
			if berr != nil {
				log.Err(berr).Msg("failed to get block number")
				wc.Close()
				return
			}
			Cache.Set("block_number", bn, 12*time.Second)
			wc.Close()
			log.Info().Msg(fmt.Sprintf("Received new timestamp: %s", t))
		}
	}
}
