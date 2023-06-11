package web3_client

import (
	"context"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

type PricingData struct {
	v2Pair UniswapV2Pair
}

func (u *UniswapClient) GetPricingData(ctx context.Context, path []accounts.Address) (*PricingData, error) {
	pair, err := u.PairToPrices(ctx, path)
	if err != nil {
		return nil, err
	}
	return &PricingData{
		v2Pair: pair,
	}, nil
}
