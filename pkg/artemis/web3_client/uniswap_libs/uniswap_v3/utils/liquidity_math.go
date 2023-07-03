package utils

import (
	"math/big"

	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/constants"
)

func AddDelta(x, y *big.Int) *big.Int {
	if y.Cmp(constants.Zero) < 0 {
		return new(big.Int).Sub(x, new(big.Int).Mul(y, constants.NegativeOne))
	} else {
		return new(big.Int).Add(x, y)
	}
}
