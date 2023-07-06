package artemis_eth_units

import (
	"math/big"

	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
)

var (
	Gwei             = big.NewInt(1e9)
	Finney           = big.NewInt(1e15)
	TenFinney        = big.NewInt(1e16)
	Ether            = big.NewInt(1e18)
	TenThousandEther = EtherMultiple(10000)

	maxUINT    = "115792089237316195423570985008687907853269984665640564039457584007913129639935"
	MaxUINT, _ = new(big.Int).SetString(maxUINT, 10)
)

func NewBigInt(amount int) *big.Int {
	return new(big.Int).SetInt64(int64(amount))
}

func EtherMultiple(multiple int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(multiple)), Ether)
}

// bal := (*hexutil.Big)(eb)

func GweiMultiple(multiple int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(multiple)), Gwei)
}

func AddBigInt(val, plus *big.Int) *big.Int {
	return new(big.Int).Add(val, plus)
}

func SubBigInt(val, minus *big.Int) *big.Int {
	return new(big.Int).Sub(val, minus)
}

func MulBigInt(x, y *big.Int) *big.Int {
	return new(big.Int).Mul(x, y)
}

func DivBigInt(x, y *big.Int) *big.Int {
	return new(big.Int).Div(x, y)
}

func NewPercentFromInts(num, den int) *core_entities.Percent {
	return core_entities.NewPercent(NewBigInt(num), NewBigInt(den))
}

func IsXGreaterThanY(x, y *big.Int) bool {
	return x.Cmp(y) > 0
}

func IsXLessThanY(x, y *big.Int) bool {
	return x.Cmp(y) < 0
}

func SetSlippage(amountOut *big.Int) *big.Int {
	slippagePerc := NewPercentFromInts(1, 10000)
	slippageAmount := FractionalAmount(amountOut, slippagePerc)
	return SubBigInt(amountOut, slippageAmount)
}

func FractionalAmount(amount *big.Int, perc *core_entities.Percent) *big.Int {
	amountOut := MulBigInt(amount, perc.Numerator)
	amountOut = DivBigInt(amountOut, perc.Denominator)
	return amountOut
}

func PercentDiff(calculated, actual *big.Int) *big.Int {
	var diff big.Int
	var absDiff big.Int
	var percentDiff big.Int
	var hundred big.Int

	// Calculate the difference
	diff.Sub(calculated, actual)

	// Get absolute difference
	absDiff.Abs(&diff)

	// Multiply by 100 for percentage
	hundred.SetInt64(100)
	percentDiff.Mul(&absDiff, &hundred)

	// Divide by actual
	percentDiff.Div(&percentDiff, actual)

	return &percentDiff
}
