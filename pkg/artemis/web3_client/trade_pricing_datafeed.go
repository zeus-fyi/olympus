package web3_client

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

type PricingData struct {
	V2Pair UniswapV2Pair
	V3Pair UniswapPoolV3
}

func (u *UniswapClient) GetV2PricingData(ctx context.Context, path []accounts.Address) (*PricingData, error) {
	pair, err := u.V2PairToPrices(ctx, path)
	if err != nil {
		return nil, err
	}
	return &PricingData{
		V2Pair: pair,
	}, nil
}

func (u *UniswapClient) GetV3PricingData(ctx context.Context, path TokenFeePath) (*PricingData, error) {
	pairV3 := UniswapPoolV3{
		Web3Actions:          u.Web3Client.Web3Actions,
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
	return &PricingData{
		V3Pair: pairV3,
	}, nil
}

/*
type PricingData struct {
	v2Pair         UniswapV2Pair
	token0EthPrice *big.Int
	token0UsdPrice *big.Int
	token1EthPrice *big.Int
	token1UsdPrice *big.Int
}
	token0EthPrice, err := price_quoter.GetETHSwapQuote(ctx, pair.Token0.String())
	if err != nil {
		log.Err(err).Msg("failed to get eth price for token0")
	}
	token0UsdPrice, err := price_quoter.GetUSDSwapQuote(ctx, pair.Token0.String())
	if err != nil {
		log.Err(err).Msg("failed to get usd price for token0")
	}
	token1EthPrice, err := price_quoter.GetETHSwapQuote(ctx, pair.Token1.String())
	if err != nil {
		log.Err(err).Msg("failed to get eth price for token1")
	}
	token1UsdPrice, err := price_quoter.GetUSDSwapQuote(ctx, pair.Token1.String())
	if err != nil {
		log.Err(err).Msg("failed to get usd price for token1")
	}
*/
