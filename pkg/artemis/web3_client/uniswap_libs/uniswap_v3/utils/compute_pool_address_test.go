package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeus-fyi/gochain/web3/accounts"
	entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/constants"
)

func TestComputePoolAddress(t *testing.T) {
	factoryAddress := accounts.HexToAddress("0x1111111111111111111111111111111111111111")
	tokenA := entities.NewToken(1, accounts.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), 18, "USDC", "USD Coin")
	tokenB := entities.NewToken(1, accounts.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"), 18, "DAI", "Dai Stablecoin")
	result, err := ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeLow, "")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, result, accounts.HexToAddress("0x90B1b09A9715CaDbFD9331b3A7652B24BfBEfD32"))

	USDC := entities.NewToken(1, accounts.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), 18, "USDC", "USD Coin")
	DAI := entities.NewToken(1, accounts.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"), 18, "DAI", "Dai Stablecoin")
	resultA, err := ComputePoolAddress(factoryAddress, USDC, DAI, constants.FeeLow, "")
	if err != nil {
		panic(err)
	}
	resultB, err := ComputePoolAddress(factoryAddress, DAI, USDC, constants.FeeLow, "")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, resultA, resultB, "should correctly compute the pool address")
}
