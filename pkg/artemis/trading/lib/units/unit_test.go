package artemis_eth_units

import (
	"fmt"
	"math/big"
	"testing"
)

func TestPercentDiff(t *testing.T) {

	valOne := NewBigIntFromStr("30600000000000000000000000000")
	valTwo := NewBigIntFromStr("30596939999999999893845539780")

	diff := PercentDiff(valOne, valTwo)
	fmt.Println("diff", diff.String())

	per := new(big.Float).SetInt(diff).Quo(new(big.Float).SetInt(diff), new(big.Float).SetInt(big.NewInt(1000000)))
	fmt.Println("per", per.String())

	meets := PercentDiffFloatComparison(valOne, valTwo, 0.01)
	fmt.Println("meets", meets)
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
