package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMostSignificantBit(t *testing.T) {
	msb, err := MostSignificantBit(big.NewInt(1))
	assert.NoError(t, err)
	assert.Equal(t, int64(0), msb)

	msb, err = MostSignificantBit(big.NewInt(2))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), msb)

	for i := int64(0); i < 255; i++ {
		msb, err = MostSignificantBit(new(big.Int).Exp(big.NewInt(2), big.NewInt(i), nil))
		assert.NoError(t, err)
		assert.Equal(t, i, msb)
	}
}
