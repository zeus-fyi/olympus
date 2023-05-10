package web3_client

import (
	"fmt"
	"math/big"
)

func (s *Web3ClientTestSuite) TestCalculateSlippage() {
	reserve0, _ := new(big.Int).SetString("8352835115956369451", 10)
	reserve1, _ := new(big.Int).SetString("16292612575348703137", 10)
	price0CumulativeLast, _ := new(big.Int).SetString("16189770433890398272", 10)
	price1CumulativeLast, _ := new(big.Int).SetString("15199761958283573464", 10)
	token0Addr, token1Addr := StringsToAddresses(PepeContractAddr, WETH9ContractAddress)
	mockPairResp := UniswapV2Pair{
		PairContractAddr:     "",
		Price0CumulativeLast: price0CumulativeLast,
		Price1CumulativeLast: price1CumulativeLast,
		KLast:                big.NewInt(0),
		Token0:               token0Addr,
		Token1:               token1Addr,
		Reserve0:             reserve0,
		Reserve1:             reserve1,
	}
	fmt.Println("mockPairResp", mockPairResp)
	//mockTrade :=

	price, err := mockPairResp.GetPriceWithBaseUnit(WETH9ContractAddress)
	s.Require().Nil(err)
	fmt.Println("weth/pepe", "price", price)
}

func (s *Web3ClientTestSuite) TestGetPairContractInfoStable() {
	uni := InitUniswapV2Client(ctx, s.MainnetWeb3User)
	pairAddr := uni.GetPairContractFromFactory(ctx, WETH9ContractAddress, LinkTokenAddr)
	pair, err := uni.GetPairContractPrices(ctx, pairAddr.String())
	s.Assert().Nil(err)
	s.Assert().NotEmpty(pair)

	fmt.Println("token0", pair.Token0.String())
	fmt.Println("token1", pair.Token1.String())
	s.Assert().Equal(LinkTokenAddr, pair.Token0.String())
	s.Assert().Equal(WETH9ContractAddress, pair.Token1.String())
	fmt.Println("kLast", pair.KLast.Uint64())
	fmt.Println("reserve0", pair.Reserve0.Uint64())
	fmt.Println("reserve1", pair.Reserve1.Uint64())
	fmt.Println("price0CumulativeLast", pair.Price0CumulativeLast.Uint64())
	fmt.Println("price1CumulativeLast", pair.Price1CumulativeLast.Uint64())

	token0Price, err := pair.GetToken0Price()
	s.Require().Nil(err)
	fmt.Println("link/weth", "token0Price", token0Price, "price0CumuluativeLast", pair.Price0CumulativeLast)

	token1Price, err := pair.GetToken1Price()
	s.Require().Nil(err)
	fmt.Println("weth/link", "token1Price", token1Price, "price1CumuluativeLast", pair.Price1CumulativeLast)

	price, err := pair.GetPriceWithBaseUnit(WETH9ContractAddress)
	s.Require().Nil(err)
	fmt.Println("weth/link", "price", price)

	price, err = pair.GetPriceWithBaseUnit(LinkTokenAddr)
	s.Require().Nil(err)
	fmt.Println("link/weth", "price", price)
}
