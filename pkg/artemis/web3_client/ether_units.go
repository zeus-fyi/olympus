package web3_client

import "math/big"

var (
	Gwei             = big.NewInt(1e9)
	Finney           = big.NewInt(1e15)
	TenFinney        = big.NewInt(1e16)
	Ether            = big.NewInt(1e18)
	TenThousandEther = new(big.Int).Mul(big.NewInt(10000), Ether)
)
