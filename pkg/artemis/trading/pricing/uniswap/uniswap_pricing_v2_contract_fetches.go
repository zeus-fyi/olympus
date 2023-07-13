package artemis_uniswap_pricing

import (
	"context"
	"fmt"
	"strings"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
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

func GetPairContractPrices(ctx context.Context, wc web3_actions.Web3Actions, p *UniswapV2Pair) error {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: p.PairContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       v2ABI,
	}
	scInfo.MethodName = getReserves
	bn, berr := artemis_trading_cache.GetLatestBlockFromCacheOrProvidedSource(ctx, wc)
	if berr != nil {
		return berr
	}
	bnst := fmt.Sprintf("%d", bn)
	sessionID := wc.GetSessionLockHeader()
	if wc.GetSessionLockHeader() != "" {
		bnst = fmt.Sprintf("%s-%s", bnst, sessionID)
	}

	tag := strings.Join([]string{fmt.Sprintf("%s", p.PairContractAddr), bnst}, "-")
	if cached, found := Cache.Get(tag); found {
		if cached == nil {
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
		log.Info().Msgf("Found cached pair %s", tag)
		pair := cached.(UniswapV2Pair)
		p.Reserve0 = pair.Reserve0
		p.Reserve1 = pair.Reserve1
		p.BlockTimestampLast = pair.BlockTimestampLast
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
