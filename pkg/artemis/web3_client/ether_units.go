package web3_client

import "math/big"

var (
	Gwei             = big.NewInt(1e9)
	Finney           = big.NewInt(1e15)
	TenFinney        = big.NewInt(1e16)
	Ether            = big.NewInt(1e18)
	TenThousandEther = EtherMultiple(10000)
)

func EtherMultiple(multiple int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(multiple)), Ether)
}

// bal := (*hexutil.Big)(eb)

func GweiMultiple(multiple int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(multiple)), Gwei)
}
