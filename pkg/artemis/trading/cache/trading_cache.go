package artemis_trading_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	redis_mev "github.com/zeus-fyi/olympus/datastores/redis/apps/mev"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
)

const (
	irisSvcBeacons = "http://iris.iris.svc.cluster.local/v2/internal/router"
)

var (
	TokenMap   map[string]artemis_autogen_bases.Erc20TokenInfo
	Cache      = cache.New(12*time.Second, 4*time.Second)
	Wc         web3_actions.Web3Actions
	WriteRedis redis_mev.MevCache
	ReadRedis  redis_mev.MevCache
)

func InitProductionRedis(ctx context.Context) {
	writeRedisOpts := redis.Options{
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Addr:         "redis-master.redis.svc.cluster.local:6379",
	}
	writer := redis.NewClient(&writeRedisOpts)
	WriteRedis = redis_mev.NewMevCache(context.Background(), writer)
	readRedisOpts := redis.Options{
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Addr:         "redis-replicas.redis.svc.cluster.local:6379",
	}
	reader := redis.NewClient(&readRedisOpts)
	ReadRedis = redis_mev.NewMevCache(context.Background(), reader)
	return
}

func InitTokenFilter(ctx context.Context) {
	_, tm, terr := artemis_mev_models.SelectERC20Tokens(ctx)
	if terr != nil {
		panic(terr)
	}
	TokenMap = tm
}

func InitWeb3Client() {
	Wc = web3_actions.NewWeb3ActionsClient(irisSvcBeacons)
	Wc.AddDefaultEthereumMainnetTableHeader()
	Wc.AddBearerToken(artemis_orchestration_auth.Bearer)
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
		return GetLatestBlock(context.Background())
	}
	if Wc.NodeURL != "" && w3SessionHeader == wcSessionHeader && len(wcSessionHeader) == 0 {
		//log.Info().Interface("w3_sessionID", w3SessionHeader).Msg("same empty session lock header, using cache")
		return GetLatestBlock(context.Background())
	}
	log.Info().Str("w3_sessionID", w3SessionHeader).Str("wc_sessionID", wcSessionHeader).Msg("different session lock header, using provided source")
	w3.Dial()
	defer w3.Close()
	w3.AddMaxBlockHeightProcedureEthJsonRpcHeader()
	bn, berr := w3.C.BlockNumber(context.Background())
	if berr != nil {
		log.Err(berr).Str("w3_sessionID", w3SessionHeader).Str("wc_sessionID", wcSessionHeader).Msg("GetLatestBlockFromCacheOrProvidedSource: failed to get block number")
		return 0, berr
	}
	return bn, nil
}

func GetLatestBlock(ctx context.Context) (uint64, error) {
	//val, ok := Cache.Get(redis_mev.LatestBlockNumberCacheKey)
	//if ok && val != nil {
	//	//log.Info().Uint64("val", val.(uint64)).Msg("got block number from cache")
	//	return val.(uint64), nil
	//}
	if ReadRedis.Client != nil {
		bn, err := ReadRedis.GetLatestBlockNumber(context.Background())
		if err == nil {
			log.Debug().Uint64("bn", bn).Msg("got block number from redis")
			Cache.Set(redis_mev.LatestBlockNumberCacheKey, bn, 6*time.Second)
			return bn, nil
		} else {
			log.Err(err).Msg("GetLatestBlock: failed to get block number from redis")
			err = nil
		}
	}
	Wc.Dial()
	defer Wc.Close()
	Wc.AddMaxBlockHeightProcedureEthJsonRpcHeader()
	bn, berr := Wc.C.BlockNumber(context.Background())
	if berr != nil {
		log.Err(berr).Msg("GetLatestBlock: failed to get block number")
		return 0, berr
	}
	if WriteRedis.Client != nil {
		err := WriteRedis.AddOrUpdateLatestBlockCache(context.Background(), bn, 12*time.Second)
		if err != nil {
			log.Err(err).Uint64("bn", bn).Msg("GetLatestBlock: failed to set block number in redis")
			err = nil
		} else {
			log.Info().Uint64("bn", bn).Msg("GetLatestBlock: set block number in redis")
		}
	}
	//log.Info().Interface("bn", bn).Msg("set block number in cache")
	//Cache.Set(redis_mev.LatestBlockNumberCacheKey, bn, 6*time.Second)
	return bn, nil
}

func SetActiveTradingBlockCache(ctx context.Context, timestampChan chan time.Time) {
	if len(artemis_orchestration_auth.Bearer) == 0 {
		panic(fmt.Errorf("bearer token is empty"))
	}

	for {
		select {
		case t := <-timestampChan:
			Wc = web3_actions.NewWeb3ActionsClient(irisSvcBeacons)
			Wc.AddDefaultEthereumMainnetTableHeader()
			Wc.AddMaxBlockHeightProcedureEthJsonRpcHeader()
			Wc.AddBearerToken(artemis_orchestration_auth.Bearer)
			Wc.Dial()
			bn, berr := Wc.C.BlockNumber(context.Background())
			if berr != nil {
				log.Err(berr).Msg("failed to get block number")
				Wc.Close()
				return
			}
			Wc.Close()
			Cache.Set(redis_mev.LatestBlockNumberCacheKey, bn, 6*time.Second)
			log.Info().Msg(fmt.Sprintf("SetActiveTradingBlockCache: Received new timestamp: %s", t))
			if WriteRedis.Client != nil {
				err := WriteRedis.AddOrUpdateLatestBlockCache(context.Background(), bn, 12*time.Second)
				if err != nil {
					log.Err(err).Str("client", WriteRedis.Client.String()).Msg("SetActiveTradingBlockCache: failed to set block number in redis")
				}
			}
		}
	}
}
