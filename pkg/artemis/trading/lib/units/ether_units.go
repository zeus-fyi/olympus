package artemis_eth_units

import (
	"math/big"

	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
)

const (
	GweiPerEth = 1000000000
)

var (
	TwoTenthGwei     = big.NewInt(2e8)
	OneTenthGwei     = big.NewInt(1e8)
	Gwei             = big.NewInt(1e9)
	Finney           = big.NewInt(1e15)
	TenFinney        = big.NewInt(1e16)
	Ether            = big.NewInt(1e18)
	TenThousandEther = EtherMultiple(10000)

	maxUINT    = "115792089237316195423570985008687907853269984665640564039457584007913129639935"
	MaxUINT, _ = new(big.Int).SetString(maxUINT, 10)
)

func NewBigIntFromUint(amount uint64) *big.Int {
	return new(big.Int).SetUint64(amount)

}
func NewBigInt(amount int) *big.Int {
	return new(big.Int).SetInt64(int64(amount))
}

func NewBigIntFromStr(amount string) *big.Int {
	val, _ := new(big.Int).SetString(amount, 10)
	return val
}

func NewBigFloatFromStr(amount string) *big.Float {
	val, ok := new(big.Float).SetString(amount)
	if !ok {
		return nil
	}
	return val
}

func EtherMultiple(multiple int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(multiple)), Ether)
}

// bal := (*hexutil.Big)(eb)

func GweiMultiple(multiple int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(multiple)), Gwei)
}

func GweiFraction(multiple int, divisor int) *big.Int {
	return DivBigInt(new(big.Int).Mul(big.NewInt(int64(multiple)), Gwei), big.NewInt(int64(divisor)))
}

func AddBigInt(val, plus *big.Int) *big.Int {
	if val == nil && plus == nil {
		return NewBigInt(0)
	}
	if val == nil {
		return plus
	}
	if plus == nil {
		return val
	}
	return new(big.Int).Add(val, plus)
}

func AddBigIntUsingInt(val *big.Int, y int) *big.Int {
	yb := NewBigInt(y)
	if val == nil {
		return yb
	}
	return new(big.Int).Add(val, yb)
}

func SubBigInt(val, minus *big.Int) *big.Int {
	if val == nil && minus == nil {
		return NewBigInt(0)
	}
	if val == nil {
		return new(big.Int).Neg(minus)
	}
	if minus == nil {
		return val
	}
	return new(big.Int).Sub(val, minus)
}

func SubUint64FBigInt(val *big.Int, uintVal uint64) *big.Int {
	return new(big.Int).Sub(val, new(big.Int).SetUint64(uintVal))
}

func MulBigInt(x, y *big.Int) *big.Int {
	if x == nil || y == nil {
		return NewBigInt(0)
	}
	return new(big.Int).Mul(x, y)
}

func MulBigIntWithInt(x *big.Int, y int) *big.Int {
	return new(big.Int).Mul(x, new(big.Int).SetInt64(int64(y)))
}

func MulBigIntWithUint64(x *big.Int, y uint64) *big.Int {
	return new(big.Int).Mul(x, new(big.Int).SetUint64(y))
}

func MulBigIntWithFloat(x *big.Int, y float64) *big.Int {
	xf := new(big.Float).SetInt(x)
	val := new(big.Float).Mul(xf, new(big.Float).SetFloat64(y))

	biReturn, _ := val.Int(nil)
	return biReturn
}

func DivBigInt(x, y *big.Int) *big.Int {
	return new(big.Int).Div(x, y)
}

func DivBigIntToFloat(x, y *big.Int) *big.Float {
	xf := new(big.Float).SetInt(x)
	xy := new(big.Float).SetInt(y)
	return new(big.Float).Quo(xf, xy)
}

func NewPercentFromInts(num, den int) *core_entities.Percent {
	return core_entities.NewPercent(NewBigInt(num), NewBigInt(den))
}

func IsXGreaterThanY(x, y *big.Int) bool {
	return x.Cmp(y) > 0
}

func IsXLessThanEqZeroOrOne(x *big.Int) bool {
	if x == nil {
		return true
	}
	if IsXLessThanY(x, NewBigInt(2)) {
		return true
	}
	return false
}

func AreAnyValuesLessThanEqZeroOrOne(x ...*big.Int) bool {
	for _, val := range x {
		if IsXLessThanEqZeroOrOne(val) {
			return true
		}
	}
	return false
}

func IsXGreaterThanZero(x *big.Int) bool {
	return x.Cmp(NewBigInt(0)) > 0
}

func IsXGreaterThanOrEqualToY(x, y *big.Int) bool {
	if x.String() == y.String() {
		return true
	}
	return IsXGreaterThanY(x, y)
}

func IsXLessThanY(x, y *big.Int) bool {
	return x.Cmp(y) < 0
}

func SetSlippage(amountOut *big.Int) *big.Int {
	//slippagePerc := NewPercentFromInts(1, 5000)
	slippagePerc := NewPercentFromInts(1, 1000)
	slippageAmount := FractionalAmount(amountOut, slippagePerc)
	return SubBigInt(amountOut, slippageAmount)
}

func ApplyTransferTax(amountOut *big.Int, num, den int) *big.Int {
	if amountOut == nil {
		return NewBigInt(0)
	}
	slippagePerc := NewPercentFromInts(num, den)
	slippageAmount := FractionalAmount(amountOut, slippagePerc)
	return SubBigInt(amountOut, slippageAmount)
}

func FractionalAmount(amount *big.Int, perc *core_entities.Percent) *big.Int {
	if amount == nil || perc == nil {
		return NewBigInt(0)
	}
	amountOut := MulBigInt(amount, perc.Numerator)
	if perc.Denominator == nil || perc.Denominator.String() == "0" {
		return NewBigInt(0)
	}
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

func PercentDiffHighPrecision(calculated, actual *big.Int) *big.Int {
	diff := new(big.Int)
	absDiff := new(big.Int)
	percentDiff := new(big.Int)
	hundred := big.NewInt(1000000) // 100*10000

	// Calculate the difference
	diff.Sub(calculated, actual)

	// Get absolute difference
	absDiff.Abs(diff)

	// Multiply by 100 for percentage and by 10000 for precision
	percentDiff.Mul(absDiff, hundred)

	// Divide by actual
	percentDiff.Div(percentDiff, actual)

	return percentDiff
}

func PercentDiffFloat(calculated, actual *big.Int) float64 {
	diff := PercentDiffHighPrecision(calculated, actual)
	percentDiff := new(big.Float).SetInt(diff).Quo(new(big.Float).SetInt(diff), new(big.Float).SetInt(big.NewInt(1000000)))
	val, _ := percentDiff.Float64()
	return val
}

func PercentDiffFloatComparison(calculated, actual *big.Int, percentCriteria float64) bool {
	diff := PercentDiffFloat(calculated, actual)
	if diff <= percentCriteria {
		return true
	}
	return false
}
