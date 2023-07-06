package artemis_uniswap_pricing

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

type UniswapPricingData struct {
	V2Pair UniswapV2Pair
	V3Pair UniswapV3Pair
}

func V2PairToPrices(ctx context.Context, wc web3_actions.Web3Actions, pairAddr []accounts.Address) (UniswapV2Pair, error) {
	p := UniswapV2Pair{}
	if len(pairAddr) == 2 {
		err := p.PairForV2(pairAddr[0].String(), pairAddr[1].String())
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: PairForV2")
			return p, err
		}
		err = GetPairContractPrices(ctx, wc, &p)
		if err != nil {
			log.Err(err).Msg("V2PairToPrices: GetPairContractPrices")
			return p, err
		}
		return p, err
	}
	return UniswapV2Pair{}, errors.New("pair address length is not 2, multi-hops not implemented yet")
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
	simMode := false
	err := pairV3.PricingData(ctx, path, simMode)
	if err != nil {
		log.Err(err).Interface("path", path).Interface("simMode", simMode).Msg("error getting v3 pricing data")
		return &UniswapPricingData{
			V3Pair: pairV3,
		}, err
	}
	return &UniswapPricingData{
		V3Pair: pairV3,
	}, nil
}
