package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeastSignificantBit(t *testing.T) {
	lsb, err := LeastSignificantBit(big.NewInt(1))
	assert.NoError(t, err)

	assert.Equal(t, int64(0), lsb)
	lsb, err = LeastSignificantBit(big.NewInt(2))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lsb)

	lsb, err = LeastSignificantBit(big.NewInt(4))
	assert.NoError(t, err)
	assert.Equal(t, int64(2), lsb)

	// Test all powers of 2 up to 255
	for i := int64(0); i < 255; i++ {
		lsb, err = LeastSignificantBit(new(big.Int).Exp(big.NewInt(2), big.NewInt(i), nil))
		assert.NoError(t, err)
		assert.Equal(t, i, lsb)
	}
}
