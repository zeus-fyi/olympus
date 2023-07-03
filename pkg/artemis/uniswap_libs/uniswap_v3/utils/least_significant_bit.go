package utils

import (
	"errors"
	"math/big"
)

func LeastSignificantBit(x *big.Int) (int64, error) {
	if x.Cmp(big.NewInt(0)) <= 0 {
		return 0, errors.New("input must be greater than 0")
	}
	y := new(big.Int).Set(x)
	i := 0
	for y.BitLen() > 0 {
		if y.Bit(0) == 1 {
			return int64(i), nil
		}
		y.Rsh(y, 1)
		i++
	}
	return int64(i), nil
}
