package artemis_eth_units

import "math/big"

var (
	Gwei             = big.NewInt(1e9)
	Finney           = big.NewInt(1e15)
	TenFinney        = big.NewInt(1e16)
	Ether            = big.NewInt(1e18)
	TenThousandEther = EtherMultiple(10000)

	maxUINT    = "115792089237316195423570985008687907853269984665640564039457584007913129639935"
	MaxUINT, _ = new(big.Int).SetString(maxUINT, 10)
)

func EtherMultiple(multiple int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(multiple)), Ether)
}

// bal := (*hexutil.Big)(eb)

func GweiMultiple(multiple int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(multiple)), Gwei)
}
