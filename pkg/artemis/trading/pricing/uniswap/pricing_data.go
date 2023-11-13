package artemis_uniswap_pricing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/patrickmn/go-cache"
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

func V2PairToPrices(ctx context.Context, bn uint64, wc web3_actions.Web3Actions, pairAddr []accounts.Address) (*UniswapV2Pair, error) {
	p := &UniswapV2Pair{}
	if len(pairAddr) == 2 {
		err := p.PairForV2(pairAddr[0].String(), pairAddr[1].String())
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: PairForV2")
			return nil, err
		}
		err = GetPairContractPrices(ctx, bn, wc, p)
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
	bn, berr := artemis_trading_cache.GetLatestBlockFromCacheOrProvidedSource(context.Background(), wc)
	if berr != nil {
		log.Err(berr).Msg("GetPairContractPrices: failed to get latest block from cache or provided source")
		return nil, berr
	}
	if len(pairAddr) == 2 {
		err := p.PairForV2(pairAddr[0].String(), pairAddr[1].String())
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: PairForV2")
			return nil, err
		}
		err = GetPairContractPrices(ctx, bn, wc, &p)
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: GetPairContractPrices")
			return nil, err
		}
		pd := &UniswapPricingData{V2Pair: p}
		return pd, err
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
	bn, berr := artemis_trading_cache.GetLatestBlockFromCacheOrProvidedSource(context.Background(), wc)
	if berr != nil {
		return nil, berr
	}
	sessionID := wc.GetSessionLockHeader()
	err = redisCache.AddV3PairToNextLookupSet(ctx, pairV3.PoolAddress, bn, sessionID)
	if err != nil {
		log.Err(err).Interface("path", path).Msg("GetV3PricingData: error adding v3 pair to next lookup set")
		err = nil
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
	statusCmd := m.Set(context.Background(), tag, bin, ttl)
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
	err := m.Get(context.Background(), tag).Scan(&bytes)
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
	V3PairNextLookupSet = "V3PairNextLookupSet"
)

func GetPairBnCacheKey(bn uint64) string {
	return fmt.Sprintf("%s-%d", V2PairNextLookupSet, bn)
}

func GetV2PairBnCacheKeyTag(bn uint64, v2Address, sessionID string) string {
	return fmt.Sprintf("%s-%d-%s-%s", V2PairNextLookupSet, bn, v2Address, sessionID)
}

func GetV3PairBnCacheKey(bn uint64) string {
	return fmt.Sprintf("%s-%d", V3PairNextLookupSet, bn)
}
func GetV3PairBnCacheKeyTag(bn uint64, v3PairAddress, sessionID string) string {
	return fmt.Sprintf("%s-%d-%s-%s", V3PairNextLookupSet, bn, v3PairAddress, sessionID)
}

var localCache = cache.New(2*time.Minute, 2*time.Minute)

func (m *PricingCache) AddV3PairToNextLookupSet(ctx context.Context, v3pairAddr string, bn uint64, sessionID string) error {
	if artemis_trading_cache.WriteRedis.Client == nil {
		return errors.New("V3PairNextLookupSet: redis client is nil")
	}
	if sessionID != "" {
		return nil
	}
	tag := GetV3PairBnCacheKey(bn)
	if _, found := localCache.Get(tag); found {
		return nil
	}
	localCache.Set(tag, true, cache.DefaultExpiration)
	m.Client = artemis_trading_cache.WriteRedis.Client
	j := 0 // used to bypass error if it's been seen once
	times := 11
	for i := 1; i < times; i++ {
		nextBlock := bn + uint64(i)
		statusCmd := m.Client.SAdd(context.Background(), GetV3PairBnCacheKey(nextBlock), v3pairAddr)
		if statusCmd.Err() != nil && j == 0 {
			log.Ctx(ctx).Err(statusCmd.Err()).Msgf("V3PairNextLookupSet: %s", v3pairAddr)
			return statusCmd.Err()
		}
		m.Client.Expire(context.Background(), GetV3PairBnCacheKey(nextBlock), time.Hour*3)
		localCache.Set(GetV3PairBnCacheKey(nextBlock), true, cache.DefaultExpiration)
		j++
	}

	// Also set an expiration time for the set if needed
	return nil
}

func (m *PricingCache) AddV2PairToNextLookupSet(ctx context.Context, bn uint64, v2pairAddr, sessionID string) error {
	if artemis_trading_cache.WriteRedis.Client == nil {
		return errors.New("AddV2PairToNextLookupSet: redis client is nil")
	}
	if sessionID != "" {
		return nil
	}
	tag := GetPairBnCacheKey(bn)
	if _, found := localCache.Get(tag); found {
		log.Info().Msgf("AddV2PairToNextLookupSet: %s already in local cache", v2pairAddr)
		return nil
	}
	localCache.Set(tag, true, cache.DefaultExpiration)
	m.Client = artemis_trading_cache.WriteRedis.Client

	j := 0 // used to bypass error if it's been seen once
	for i := 1; i < 10; i++ {
		nextBlock := bn + 1
		statusCmd := m.Client.SAdd(context.Background(), GetPairBnCacheKey(nextBlock), v2pairAddr)
		if statusCmd.Err() != nil && j == 0 {
			log.Ctx(ctx).Err(statusCmd.Err()).Msgf("AddV2PairToNextLookupSet: %s", v2pairAddr)
			return statusCmd.Err()
		}
		m.Client.Expire(ctx, GetPairBnCacheKey(nextBlock), time.Hour*12)
		localCache.Set(GetPairBnCacheKey(nextBlock), true, cache.DefaultExpiration)
		j++
	}
	m.Client.Expire(ctx, tag, time.Hour*12)
	// Also set an expiration time for the set if needed
	return nil
}

func FetchV2PairsToMulticall(ctx context.Context, wc web3_actions.Web3Actions, bn uint64) error {
	if artemis_trading_cache.ReadRedis.Client == nil {
		return errors.New("FetchV2PairsToMulticall: redis client is nil")
	}
	redisCache.Client = artemis_trading_cache.ReadRedis.Client
	addresses, err := redisCache.GetV2PairsToMulticall(context.Background(), bn)
	if err != nil {
		return err
	}
	var tmp []string
	for _, addr := range addresses {
		tmp = append(tmp, addr)
		if len(tmp) >= 25 {
			_, err = GetBatchPairContractPricesViaMulticall3(context.Background(), wc, tmp...)
			if err != nil {
				return err
			}
			tmp = []string{}
		}
	}
	if len(tmp) > 0 {
		_, err = GetBatchPairContractPricesViaMulticall3(context.Background(), wc, tmp...)
		if err != nil {
			return err
		}
	}
	log.Info().Int("pairCount", len(addresses)).Msg("Fetched V2 pairs to multicall")
	return nil
}

func (m *PricingCache) GetV2PairsToMulticall(ctx context.Context, bn uint64) ([]string, error) {
	if artemis_trading_cache.ReadRedis.Client == nil {
		return nil, errors.New("GetV2PairsToMulticall: redis client is nil")
	}
	m.Client = artemis_trading_cache.ReadRedis.Client
	pairAddresses, err := m.Client.SMembers(context.Background(), GetPairBnCacheKey(bn)).Result()
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("GetV2PairsToMulticall: %s", GetPairBnCacheKey(bn))
		return nil, err
	}

	return pairAddresses, nil
}
