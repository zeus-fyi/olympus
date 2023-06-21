package utils

import (
	"fmt"
	"math/big"
)

func LeastSignificantBit(x *big.Int) (*big.Int, error) {
	zero := big.NewInt(0)
	if x.Cmp(zero) <= 0 {
		return nil, fmt.Errorf("input must be greater than 0")
	}

	// 1 in big.Int
	one := big.NewInt(1)
	// Initialize the result
	r := big.NewInt(255)
	for i := 0; i <= 255; i++ {
		// Check if the bit at position i is set.
		// If it is, we have found the least significant bit
		if new(big.Int).And(x, one).Cmp(one) == 0 {
			break
		}
		// Shift the bits to the right
		x.Rsh(x, 1)
		// Decrement the result
		r.Sub(r, one)
	}
	return r, nil
}
