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
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
)

const (
	irisSvcBeacons = "http://iris.iris.svc.cluster.local/v1beta/internal/router/group?routeGroup=quiknode-mainnet"
)

var (
	TokenMap map[string]artemis_autogen_bases.Erc20TokenInfo
	Cache    = cache.New(12*time.Second, 4*time.Second)
	Wc       web3_actions.Web3Actions
)

func InitTokenFilter(ctx context.Context) {
	_, tm, terr := artemis_mev_models.SelectERC20Tokens(ctx)
	if terr != nil {
		panic(terr)
	}
	TokenMap = tm
}

func InitWeb3Client() {
	Wc = web3_actions.NewWeb3ActionsClient(irisSvcBeacons)
	if len(artemis_orchestration_auth.Bearer) == 0 {
		panic(fmt.Errorf("bearer token is empty"))
	}
	Wc.AddBearerToken(artemis_orchestration_auth.Bearer)
}

func GetLatestBlockFromCacheOrProvidedSource(ctx context.Context, w3 web3_actions.Web3Actions) (uint64, error) {
	w3SessionHeader := w3.GetSessionLockHeader()
	wcSessionHeader := Wc.GetSessionLockHeader()
	if Wc.NodeURL != "" && len(wcSessionHeader) > 0 && len(w3SessionHeader) > 0 && w3SessionHeader == wcSessionHeader {
		//log.Info().Interface("w3_sessionID", w3SessionHeader).Msg("same session lock header, using cache")
		return GetLatestBlock(ctx)
	}
	if Wc.NodeURL != "" && w3SessionHeader == wcSessionHeader && len(wcSessionHeader) == 0 {
		//log.Info().Interface("w3_sessionID", w3SessionHeader).Msg("same empty session lock header, using cache")
		return GetLatestBlock(ctx)
	}
	log.Info().Str("w3_sessionID", w3SessionHeader).Str("wc_sessionID", wcSessionHeader).Msg("different session lock header, using provided source")
	w3.Dial()
	defer w3.Close()
	bn, berr := w3.C.BlockNumber(ctx)
	if berr != nil {
		log.Err(berr).Msg("failed to get block number")
		return 0, berr
	}
	return bn, berr
}

func GetLatestBlock(ctx context.Context) (uint64, error) {
	val, ok := Cache.Get("block_number")
	if ok && val != nil {
		//log.Info().Uint64("val", val.(uint64)).Msg("got block number from cache")
		return val.(uint64), nil
	}
	Wc.Dial()
	defer Wc.Close()
	bn, berr := Wc.C.BlockNumber(ctx)
	if berr != nil {
		log.Err(berr).Msg("failed to get block number")
		return 0, berr
	}
	//log.Info().Interface("bn", bn).Msg("set block number in cache")
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
			Wc = web3_actions.NewWeb3ActionsClient(irisSvcBeacons)
			Wc.AddBearerToken(artemis_orchestration_auth.Bearer)
			Wc.Dial()
			bn, berr := Wc.C.BlockNumber(ctx)
			if berr != nil {
				log.Err(berr).Msg("failed to get block number")
				Wc.Close()
				return
			}
			Cache.Set("block_number", bn, 12*time.Second)
			Wc.Close()
			log.Info().Msg(fmt.Sprintf("Received new timestamp: %s", t))
		}
	}
}
