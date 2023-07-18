package artemis_uniswap_pricing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

var redisCache = PricingCache{nil}

type UniswapPricingData struct {
	V2Pair UniswapV2Pair
	V3Pair UniswapV3Pair
}

func V2PairToPrices(ctx context.Context, wc web3_actions.Web3Actions, pairAddr []accounts.Address) (*UniswapV2Pair, error) {
	p := &UniswapV2Pair{}
	if len(pairAddr) == 2 {
		err := p.PairForV2(pairAddr[0].String(), pairAddr[1].String())
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: PairForV2")
			return nil, err
		}
		err = GetPairContractPrices(ctx, wc, p)
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: GetPairContractPrices")
			return nil, err
		}
		return p, err
	}
	return nil, errors.New("pair address length is not 2, multi-hops not implemented yet")
}

//func GetLiveV2PricingData(ctx context.Context, pairAddr []accounts.Address) (*UniswapPricingData, error) {
//	wc := web3_actions.Web3Actions{}
//	return GetV2PricingData(ctx, wc, pairAddr)
//}

func GetV2PricingData(ctx context.Context, wc web3_actions.Web3Actions, pairAddr []accounts.Address) (*UniswapPricingData, error) {
	p := UniswapV2Pair{}
	if len(pairAddr) == 2 {
		err := p.PairForV2(pairAddr[0].String(), pairAddr[1].String())
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: PairForV2")
			return nil, err
		}
		err = GetPairContractPrices(ctx, wc, &p)
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: GetPairContractPrices")
			return nil, err
		}
		return &UniswapPricingData{V2Pair: p}, err
	}
	return nil, errors.New("pair address length is not 2, multi-hops not implemented yet")
}

func GetV3PricingData(ctx context.Context, wc web3_actions.Web3Actions, path artemis_trading_types.TokenFeePath) (*UniswapPricingData, error) {
	pairV3 := UniswapV3Pair{
		Web3Actions:          wc,
		PoolAddress:          "",
		Fee:                  0,
		Slot0:                Slot0{},
		Liquidity:            nil,
		TickListDataProvider: nil,
	}
	err := pairV3.PricingData(ctx, path)
	if err != nil {
		log.Err(err).Interface("path", path).Msg("error getting v3 pricing data")
		return nil, err
	}
	return &UniswapPricingData{
		V3Pair: pairV3,
	}, nil
}

type PricingCache struct {
	*redis.Client
}

func (m *PricingCache) AddOrUpdatePairPricesCache(ctx context.Context, tag string, pd UniswapPricingData, ttl time.Duration) error {
	if artemis_trading_cache.WriteRedis.Client == nil {
		return errors.New("AddOrUpdatePairPricesCache: redis client is nil")
	}
	m.Client = artemis_trading_cache.WriteRedis.Client
	bin, err := m.MarshalBinary(pd)
	if err != nil {
		return err
	}
	statusCmd := m.Set(ctx, tag, bin, ttl)
	if statusCmd.Err() != nil {
		log.Ctx(ctx).Err(statusCmd.Err()).Msgf("AddOrUpdateLatestBlockCache: %s", tag)
		return statusCmd.Err()
	}
	return nil
}

func (m *PricingCache) MarshalBinary(up UniswapPricingData) ([]byte, error) {
	return json.Marshal(up)
}

func (m *PricingCache) UnmarshalBinary(data []byte) (UniswapPricingData, error) {
	pd := UniswapPricingData{}
	err := json.Unmarshal(data, &pd)
	return pd, err
}

func (m *PricingCache) GetPairPricesFromCacheIfExists(ctx context.Context, tag string) (UniswapPricingData, error) {
	if artemis_trading_cache.ReadRedis.Client == nil {
		return UniswapPricingData{}, errors.New("AddOrUpdatePairPricesCache: redis client is nil")
	}
	m.Client = artemis_trading_cache.ReadRedis.Client
	pd := UniswapPricingData{}
	var bytes []byte
	err := m.Get(ctx, tag).Scan(&bytes)
	switch {
	case err == redis.Nil:
		return pd, fmt.Errorf("GetPairPricesFromCacheIfExists: %s", tag)
	case err != nil:
		log.Err(err).Msgf("GetPairPricesFromCacheIfExists Get failed: %s", tag)
	}
	cachedPd, err := m.UnmarshalBinary(bytes)
	if err != nil {
		return pd, err
	}
	return cachedPd, nil
}

const (
	V2PairNextLookupSet = "V2PairNextLookupSet"
)

func GetPairBnCacheKey(bn uint64) string {
	return fmt.Sprintf("%s-%d", V2PairNextLookupSet, bn)
}

func (m *PricingCache) AddV2PairToNextLookupSet(ctx context.Context, v2pairAddr string, bn uint64) error {
	if artemis_trading_cache.WriteRedis.Client == nil {
		return errors.New("AddV2PairToNextLookupSet: redis client is nil")
	}
	m.Client = artemis_trading_cache.WriteRedis.Client
	statusCmd := m.Client.SAdd(ctx, GetPairBnCacheKey(bn), v2pairAddr)
	if statusCmd.Err() != nil {
		log.Ctx(ctx).Err(statusCmd.Err()).Msgf("AddV2PairToNextLookupSet: %s", v2pairAddr)
		return statusCmd.Err()
	}
	// Also set an expiration time for the set if needed
	m.Client.Expire(ctx, GetPairBnCacheKey(bn), time.Hour*12)
	return nil
}

func (m *PricingCache) GetV2PairsToMulticall(ctx context.Context, bn uint64) ([]string, error) {
	if artemis_trading_cache.ReadRedis.Client == nil {
		return nil, errors.New("GetV2PairsToMulticall: redis client is nil")
	}
	m.Client = artemis_trading_cache.ReadRedis.Client
	pairAddresses, err := m.Client.SMembers(ctx, GetPairBnCacheKey(bn)).Result()
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("GetV2PairsToMulticall: %s", GetPairBnCacheKey(bn))
		return nil, err
	}
	return pairAddresses, nil
}
