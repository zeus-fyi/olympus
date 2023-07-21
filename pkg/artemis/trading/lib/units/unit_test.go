package artemis_eth_units

import (
	"fmt"
	"math/big"
	"testing"
)

func TestPercentDiff(t *testing.T) {

	valOne := NewBigIntFromStr("29999643121717752")
	valTwo := NewBigIntFromStr("29996643157405581")

	diff := PercentDiff(valOne, valTwo)
	fmt.Println("diff", diff.String())

	per := PercentDiffFloat(valOne, valTwo)
	fmt.Println("per", per)

	meets := PercentDiffFloatComparison(valOne, valTwo, 0.0001)
	fmt.Println("meets", meets)

	f := NewPercentFromInts(0, 0)
	fmt.Println("f", f.Numerator)

	startVal := NewBigIntFromStr("0")
	fmt.Println("startVal", startVal.String())
	endVal := ApplyTransferTax(startVal, 1, 1)
	fmt.Println("val", endVal.String())

	resp := MulBigInt(NewBigInt(0), NewBigInt(0))
	fmt.Println("resp", resp.String())

	bo0 := IsXLessThanEqZeroOrOne(NewBigInt(0))
	fmt.Println("bo", bo0)
	bo1 := IsXLessThanEqZeroOrOne(NewBigInt(1))
	fmt.Println("bo", bo1)
	bo2 := IsXLessThanEqZeroOrOne(NewBigInt(2))
	fmt.Println("bo", bo2)
}

func TestPercentDiffk(t *testing.T) {
	tests := []struct {
		name string
		calc *big.Int
		act  *big.Int
		want *big.Int
	}{
		{
			name: "Case 1: 100 vs 50",
			calc: big.NewInt(100),
			act:  big.NewInt(50),
			want: big.NewInt(1000000), // 100.0000 percent difference
		},
		{
			name: "Case 2: 200 vs 150",
			calc: big.NewInt(200),
			act:  big.NewInt(150),
			want: big.NewInt(333333), // 33.3333 percent difference
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PercentDiff(tt.calc, tt.act); got.Cmp(tt.want) != 0 {
				t.Errorf("PercentDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}
