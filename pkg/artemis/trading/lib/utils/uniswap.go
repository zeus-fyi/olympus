package artemis_utils

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
)

const (
	pairAddressSuffix = "96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f"
)

func CreateV2TradingPair(tokenA, tokenB accounts.Address) (string, accounts.Address, accounts.Address) {
	if tokenA.Str() == artemis_trading_constants.ZeroAddress {
		tokenA = artemis_trading_constants.WETH9ContractAddressAccount
	}
	if tokenB.Str() == artemis_trading_constants.ZeroAddress {
		tokenB = artemis_trading_constants.WETH9ContractAddressAccount
	}
	if tokenA == tokenB {
		panic("tokenA and tokenB cannot be the same")
	}
	token0, token1 := SortTokens(tokenA, tokenB)
	message := []byte{255}
	message = append(message, common.HexToAddress(artemis_trading_constants.UniswapV2FactoryAddress).Bytes()...)
	addrSum := token0.Bytes()
	addrSum = append(addrSum, token1.Bytes()...)
	message = append(message, crypto.Keccak256(addrSum)...)
	b, err := hex.DecodeString(pairAddressSuffix)
	if err != nil {
		panic(err)
	}
	message = append(message, b...)
	hashed := crypto.Keccak256(message)
	addressBytes := big.NewInt(0).SetBytes(hashed)
	addressBytes = addressBytes.Abs(addressBytes)
	pairAddr := common.BytesToAddress(addressBytes.Bytes()).String()
	return pairAddr, token0, token1
}

/*
	pair, token0, token1 := artemis_utils.CreateV2TradingPair(artemis_trading_constants.PepeContractAddr, artemis_trading_constants.WETH9ContractAddress)
	s.Assert().NotEmpty(pair)
	fmt.Println(pair)
	fmt.Println(token0.String())
	fmt.Println(token1.String())

	pair2, token0a, token1a := artemis_utils.CreateV2TradingPair(artemis_trading_constants.WETH9ContractAddress, artemis_trading_constants.PepeContractAddr)
	s.Assert().NotEmpty(pair)

	s.Assert().Equal(token0.String(), token0a.String())
	s.Assert().Equal(token1.String(), token1a.String())
	s.Assert().Equal(pair, pair2)
*/
