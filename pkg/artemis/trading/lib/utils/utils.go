package artemis_utils

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

func DivideByHalf(input *big.Int) *big.Int {
	modEven := new(big.Int).Mod(input, big.NewInt(2))
	if modEven.String() == "0" {
		input = input.Div(input, big.NewInt(2))
	} else {
		input = input.Add(input, big.NewInt(1))
		input = input.Div(input, big.NewInt(2))
	}
	return input
}

func SortTokens(tkn0, tkn1 accounts.Address) (accounts.Address, accounts.Address) {
	token0Rep := big.NewInt(0).SetBytes(tkn0.Bytes())
	token1Rep := big.NewInt(0).SetBytes(tkn1.Bytes())

	if token0Rep.Cmp(token1Rep) > 0 {
		tkn0, tkn1 = tkn1, tkn0
	}
	return tkn0, tkn1
}
