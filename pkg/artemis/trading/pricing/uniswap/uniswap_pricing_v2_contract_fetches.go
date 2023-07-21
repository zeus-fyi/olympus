package artemis_uniswap_pricing

import (
	"context"
	"fmt"
	"strings"

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
	bnst := fmt.Sprintf("%d", bn)
	sessionID := wc.GetSessionLockHeader()
	if wc.GetSessionLockHeader() != "" {
		bnst = fmt.Sprintf("%s-%s", bnst, sessionID)
	}
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
		if p != nil {
			tag := strings.Join([]string{fmt.Sprintf("%s", p.PairContractAddr), bnst}, "-")
			Cache.Set(tag, *p, cache.NoExpiration)
		}
	}
	return pairs, nil
}

func GetPairContractPrices(ctx context.Context, wc web3_actions.Web3Actions, p *UniswapV2Pair) error {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: p.PairContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       v2ABI,
	}
	scInfo.MethodName = getReserves
	bn, berr := artemis_trading_cache.GetLatestBlockFromCacheOrProvidedSource(context.Background(), wc)
	if berr != nil {
		log.Err(berr).Msg("GetPairContractPrices: failed to get latest block from cache or provided source")
		return berr
	}
	bnst := fmt.Sprintf("%d", bn)
	sessionID := wc.GetSessionLockHeader()
	if wc.GetSessionLockHeader() != "" {
		bnst = fmt.Sprintf("%s-%s", bnst, sessionID)
	} else {
		err := redisCache.AddV2PairToNextLookupSet(context.Background(), bn, p.PairContractAddr, sessionID)
		if err != nil {
			log.Error().Err(err).Msg("AddV2PairToNextLookupSet: failed to add pair to next lookup set")
			err = nil
		}
	}
	tag := strings.Join([]string{fmt.Sprintf("%s", p.PairContractAddr), bnst}, "-")
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
		if cached == nil {
			resp, err := wc.CallConstantFunction(ctx, scInfo)
			if err != nil {
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
			Cache.Set(tag, *p, cache.NoExpiration)
			return nil
		}
		if sessionID != "" {
			log.Info().Msgf("Found cached pair %s", tag)
		}
		pair := cached.(UniswapV2Pair)
		p.Reserve0 = pair.Reserve0
		p.Reserve1 = pair.Reserve1
		p.BlockTimestampLast = pair.BlockTimestampLast
		return nil
	}
	pd, err := redisCache.GetPairPricesFromCacheIfExists(context.Background(), tag)
	if err != nil {
		log.Err(err).Msgf("GetPairContractPrices: GetPairPricesFromCacheIfExists: error getting pair prices from cache for %s", tag)
		err = nil
	} else {
		cachedV2 := pd.V2Pair
		p = &cachedV2
		return nil
	}
	resp, err := wc.CallConstantFunction(ctx, scInfo)
	if err != nil {
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
	Cache.Set(tag, *p, cache.NoExpiration)
	if len(resp) <= 2 {
		return err
	}
	return nil
}
