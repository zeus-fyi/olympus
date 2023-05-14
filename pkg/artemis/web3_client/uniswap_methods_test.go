package web3_client

import (
	"fmt"
	"math/big"
)

// https://www.defi-sandwi.ch/

// upper-bound on the profit is the victim’s trade amount
// sandwich attacker needs to pay the 0.3% fee twice
func (s *Web3ClientTestSuite) TestSandwichAttack() {
	/*
		Swaps an exact amount of tokens for as much ETH as possible, along the route determined by the path. The first element of path is the input token,
		the last must be WETH, and any intermediate elements represent intermediate pairs to trade through (if, for example, a direct pair does not exist).
	*/
	amountIn, _ := new(big.Int).SetString("100000000000000000000", 10)
	amountOut, _ := new(big.Int).SetString("3223835795348941600", 10)

	// 1% slippage, meaning they're willing to receive 1% less than the amountOut as minimum acceptable amount
	slippage := new(big.Int).Div(amountOut, big.NewInt(100))
	fmt.Println("slippage", slippage.String())
	amountOutMin := new(big.Int).Sub(amountOut, slippage)
	slippageMargin := new(big.Int).Div(amountOut, big.NewInt(10000))
	amountOutMinWithMargin := new(big.Int).Add(amountOutMin, slippageMargin)
	fmt.Println("amountOutMin", amountOutMin)
	mockTrade := SwapExactTokensForETHParams{
		AmountIn:     amountIn,
		AmountOutMin: amountOut,
	}

	reserve0, _ := new(big.Int).SetString("47956013761392256000", 10)
	reserve1, _ := new(big.Int).SetString("1383382537550055000000", 10)
	token0Addr, token1Addr := StringsToAddresses(PepeContractAddr, WETH9ContractAddress)
	mockPairResp := UniswapV2Pair{
		KLast:    big.NewInt(0),
		Token0:   token0Addr,
		Token1:   token1Addr,
		Reserve0: reserve0,
		Reserve1: reserve1,
	}
	var profitString []string
	var frontRunAmounts []string
	startOffset := big.NewInt(0)
	for true {
		mockPairResp = UniswapV2Pair{
			KLast:    big.NewInt(0),
			Token0:   token0Addr,
			Token1:   token1Addr,
			Reserve0: reserve0,
			Reserve1: reserve1,
		}
		fmt.Println("-----------front run trade-----------")
		tokenSellAmount, _ := new(big.Int).SetString("3000000000000000000", 10)
		tokenSellAmount = tokenSellAmount.Add(startOffset, tokenSellAmount)
		fmt.Println("startAmount", tokenSellAmount.String())

		toFrontRun, _, _ := mockPairResp.PriceImpactToken1BuyToken0(tokenSellAmount)
		fmt.Println("endAmount", toFrontRun.AmountOut.String())
		fmt.Println("-----------user trade-----------")
		// now let user sell their tokens
		to, _, _ := mockPairResp.PriceImpactToken1BuyToken0(mockTrade.AmountIn)
		fmt.Println("userEndAmount", to.AmountOut.String())
		difference := new(big.Int).Sub(to.AmountOut, amountOutMinWithMargin)
		fmt.Println("difference", difference.String())
		if difference.Cmp(big.NewInt(0)) < 0 {
			fmt.Println("user trade failed")
			break
		}
		fmt.Println("-----------sandwich trade-----------")
		sandwichDump := toFrontRun.AmountOut
		fmt.Println("sandwichAmountToDump", sandwichDump)
		toSandwich, _, _ := mockPairResp.PriceImpactToken0BuyToken1(sandwichDump)
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		fmt.Println("endTokenAmount", toSandwich.AmountOut.String())
		fmt.Println("endProfit", profit.String())
		profitString = append(profitString, profit.String())
		frontRunAmounts = append(frontRunAmounts, toFrontRun.AmountIn.String())
		oneTenthToken, _ := new(big.Int).SetString("100000000000000000", 10)
		startOffset = new(big.Int).Add(startOffset, oneTenthToken)
	}
	fmt.Println("frontRunAmounts", frontRunAmounts)
	fmt.Println("profitAmounts", profitString)
}

func (s *Web3ClientTestSuite) TestCalculateSlippage() {
	reserve0, _ := new(big.Int).SetString("2859456211217791841082775702235", 10)
	reserve1, _ := new(big.Int).SetString("3057340484928582107066", 10)
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
	price, err := mockPairResp.GetPriceWithBaseUnit(WETH9ContractAddress)
	s.Require().Nil(err)
	fmt.Println("weth/pepe", "price", price)

	price, err = mockPairResp.GetPriceWithBaseUnit(PepeContractAddr)
	s.Require().Nil(err)
	fmt.Println("pepe/weth", "price", price)
}

func (s *Web3ClientTestSuite) TestGetPepeWETH() {
	uni := InitUniswapV2Client(ctx, s.MainnetWeb3User)
	pairAddr := uni.GetPairContractFromFactory(ctx, WETH9ContractAddress, PepeContractAddr)
	pair, err := uni.GetPairContractPrices(ctx, pairAddr.String())
	fmt.Println("token0", pair.Token0.String())
	fmt.Println("token1", pair.Token1.String())
	s.Assert().Equal(PepeContractAddr, pair.Token0.String())
	s.Assert().Equal(WETH9ContractAddress, pair.Token1.String())
	fmt.Println("kLast", pair.KLast.String())
	fmt.Println("reserve0", pair.Reserve0.String())
	fmt.Println("reserve1", pair.Reserve1.String())
	fmt.Println("price0CumulativeLast", pair.Price0CumulativeLast.Uint64())
	fmt.Println("price1CumulativeLast", pair.Price1CumulativeLast.Uint64())

	price, err := pair.GetPriceWithBaseUnit(WETH9ContractAddress)
	s.Require().Nil(err)
	fmt.Println("weth/pepe", "price", price)

	price, err = pair.GetPriceWithBaseUnit(PepeContractAddr)
	s.Require().Nil(err)
	fmt.Println("pepe/weth", "price", price)
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
	fmt.Println("reserve0", pair.Reserve0.String())
	fmt.Println("reserve1", pair.Reserve1.String())
	fmt.Println("price0CumulativeLast", pair.Price0CumulativeLast.String())
	fmt.Println("price1CumulativeLast", pair.Price1CumulativeLast.String())

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

func (s *Web3ClientTestSuite) TestGetPairContractInfoMismatchedDecimals() {
	uni := InitUniswapV2Client(ctx, s.MainnetWeb3User)
	pairAddr := uni.GetPairContractFromFactory(ctx, WETH9ContractAddress, HexTokenAddr)
	pair, err := uni.GetPairContractPrices(ctx, pairAddr.String())
	s.Assert().Nil(err)
	s.Assert().NotEmpty(pair)
	fmt.Println("token0", pair.Token0.String())
	fmt.Println("token1", pair.Token1.String())
	s.Assert().Equal(HexTokenAddr, pair.Token0.String())
	s.Assert().Equal(WETH9ContractAddress, pair.Token1.String())
	fmt.Println("kLast", pair.KLast.String())
	fmt.Println("reserve0", pair.Reserve0.String())
	fmt.Println("reserve1", pair.Reserve1.String())
	fmt.Println("price0CumulativeLast", pair.Price0CumulativeLast.String())
	fmt.Println("price1CumulativeLast", pair.Price1CumulativeLast.String())

	token0Price, err := pair.GetToken0Price()
	s.Require().Nil(err)
	fmt.Println("hex/weth", "token0Price", token0Price, "price0CumuluativeLast", pair.Price0CumulativeLast)

	token1Price, err := pair.GetToken1Price()
	s.Require().Nil(err)
	fmt.Println("weth/hex", "token1Price", token1Price, "price1CumuluativeLast", pair.Price1CumulativeLast)

	price, err := pair.GetPriceWithBaseUnit(WETH9ContractAddress)
	s.Require().Nil(err)
	fmt.Println("weth/hex", "price", price)

	price, err = pair.GetPriceWithBaseUnit(HexTokenAddr)
	s.Require().Nil(err)
	fmt.Println("hex/weth", "price", price)
}
