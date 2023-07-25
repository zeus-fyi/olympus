package artemis_uniswap_pricing

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_multicall "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/multicall"
	artemis_utils "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/utils"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

const (
	getReserves = "getReserves"
)

var (
	v2ABI = artemis_oly_contract_abis.MustLoadUniswapV2PairAbi()
	Cache = cache.New(cache.NoExpiration, cache.NoExpiration)
	// blockNumber + pairAddr -> UniswapV2Pair
)

func GetBatchPairContractPricesViaMulticall3(ctx context.Context, wc web3_actions.Web3Actions, pairAddresses ...string) ([]*UniswapV2Pair, error) {
	bn, berr := artemis_trading_cache.GetLatestBlockFromCacheOrProvidedSource(context.Background(), wc)
	if berr != nil {
		return nil, berr
	}
	sessionID := wc.GetSessionLockHeader()
	mcalls := make([]artemis_multicall.MultiCallElement, len(pairAddresses))
	for i, pairAddr := range pairAddresses {
		addr := common.HexToAddress(pairAddr)
		mcalls[i] = artemis_multicall.MultiCallElement{
			Name: getReserves,
			Call: artemis_multicall.Call{
				Target:       addr,
				AllowFailure: false,
				Data:         nil,
			},
			AbiFile:       v2ABI,
			DecodedInputs: []interface{}{},
		}
	}
	m := artemis_multicall.Multicall3{
		Calls:   mcalls,
		Results: nil,
	}
	resp, err := m.PackAndCall(ctx, wc)
	if err != nil {
		return nil, err
	}
	pairs := make([]*UniswapV2Pair, len(resp))
	for i, respVal := range resp {
		respSlice := respVal.DecodedReturnData
		p := &UniswapV2Pair{
			PairContractAddr: pairAddresses[i],
		}
		reserve0, rerr := artemis_utils.ParseBigInt(respSlice[0])
		if rerr != nil {
			return nil, rerr
		}
		p.Reserve0 = reserve0
		reserve1, rerr := artemis_utils.ParseBigInt(respSlice[1])
		if rerr != nil {
			return nil, rerr
		}
		p.Reserve1 = reserve1
		blockTimestampLast, rerr := artemis_utils.ParseBigInt(respSlice[2])
		if rerr != nil {
			return nil, rerr
		}
		p.BlockTimestampLast = blockTimestampLast
		pairs[i] = p
		if reserve0 != nil && reserve1 != nil && p != nil {
			tag := GetV2PairBnCacheKeyTag(bn, p.PairContractAddr, sessionID)
			pd := UniswapPricingData{
				V2Pair: *p,
			}
			Cache.Set(tag, *p, cache.NoExpiration)
			err = redisCache.AddOrUpdatePairPricesCache(context.Background(), tag, pd, time.Hour*24)
			if err != nil {
				log.Err(err).Msgf("GetBatchPairContractPricesViaMulticall3: error adding or updating pair prices cache for %s", tag)
				err = nil
			}
		}

	}
	return pairs, nil
}

func GetPairContractPrices(ctx context.Context, bn uint64, wc web3_actions.Web3Actions, p *UniswapV2Pair) error {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: p.PairContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       v2ABI,
	}
	scInfo.MethodName = getReserves
	sessionID := wc.GetSessionLockHeader()
	tag := GetV2PairBnCacheKeyTag(bn, p.PairContractAddr, sessionID)
	if wc.GetSessionLockHeader() == "" {
		err := redisCache.AddV2PairToNextLookupSet(context.Background(), bn, p.PairContractAddr, sessionID)
		if err != nil {
			log.Error().Err(err).Msg("AddV2PairToNextLookupSet: failed to add pair to next lookup set")
			err = nil
		}
	}
	if cached, found := Cache.Get(tag); found {
		if cached == nil {
			pd, err := redisCache.GetPairPricesFromCacheIfExists(context.Background(), tag)
			if err != nil {
				log.Err(err).Msgf("AddV2PairToNextLookupSet: error getting pair prices from cache for %s", tag)
				err = nil
			} else {
				cachedV2 := pd.V2Pair
				p = &cachedV2
				return nil
			}
		}
	}
	resp, err := wc.CallConstantFunction(ctx, scInfo)
	if err != nil {
		log.Err(err).Msgf("GetPairContractPrices: CallConstantFunction: error calling constant function for %s", tag)
		return err
	}
	if len(resp) <= 2 {
		log.Warn().Msgf("GetPairContractPrices: len(resp) <= 2 for %s", tag)
		return err
	}
	reserve0, err := artemis_utils.ParseBigInt(resp[0])
	if err != nil {
		return err
	}
	p.Reserve0 = reserve0
	reserve1, err := artemis_utils.ParseBigInt(resp[1])
	if err != nil {
		return err
	}
	p.Reserve1 = reserve1
	blockTimestampLast, err := artemis_utils.ParseBigInt(resp[2])
	if err != nil {
		return err
	}
	p.BlockTimestampLast = blockTimestampLast
	if p != nil && p.Reserve0 != nil && p.Reserve1 != nil {
		pd := UniswapPricingData{
			V2Pair: *p,
		}
		Cache.Set(tag, *p, cache.NoExpiration)
		err = redisCache.AddOrUpdatePairPricesCache(context.Background(), tag, pd, time.Hour*24)
		if err != nil {
			log.Err(err).Msgf("GetPairContractPrices: error adding or updating pair prices cache for %s", tag)
			err = nil
		}
	}
	return nil
}
