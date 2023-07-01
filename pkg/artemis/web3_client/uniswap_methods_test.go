package web3_client

import (
	"fmt"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
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
		Path: []accounts.Address{
			accounts.HexToAddress(PepeContractAddr),
			accounts.HexToAddress(WETH9ContractAddress),
		},
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
	startOffset := big.NewInt(0)
	endProfit := big.NewInt(0)
	tokenSellAmountFinal := big.NewInt(0)
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
		tokenSellAmountFinal = tokenSellAmount
		fmt.Println("-----------sandwich trade-----------")
		sandwichDump := toFrontRun.AmountOut
		fmt.Println("sandwichAmountToDump", sandwichDump)
		toSandwich, _, _ := mockPairResp.PriceImpactToken0BuyToken1(sandwichDump)
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		fmt.Println("endTokenAmount", toSandwich.AmountOut.String())
		fmt.Println("endProfit", profit.String())
		oneTenthToken, _ := new(big.Int).SetString("100000000000000000", 10)
		startOffset = new(big.Int).Add(startOffset, oneTenthToken)
		endProfit = profit
	}
	fmt.Println("--------------summary-----------")
	fmt.Println("tokenSellAmountFinal", tokenSellAmountFinal.String())
	fmt.Println("endProfit", endProfit.String())
}

func (s *Web3ClientTestSuite) TestSandwichAttackBinSearch() {
	amountIn, _ := new(big.Int).SetString("100000000000000000000", 10)
	amountOut, _ := new(big.Int).SetString("3191917650004962400", 10)

	// 1% slippage, meaning they're willing to receive 1% less than the amountOut as minimum acceptable amount
	slippage := new(big.Int).Div(amountOut, big.NewInt(100))
	fmt.Println("slippage", slippage.String())
	amountOutMin := new(big.Int).Sub(amountOut, slippage)
	//slippageMargin := new(big.Int).Div(amountOut, big.NewInt(10000))
	//amountOutMinWithMargin := new(big.Int).Add(amountOutMin, slippageMargin)
	fmt.Println("amountOutMin", amountOutMin)
	mockTrade := SwapExactTokensForETHParams{
		AmountIn:     amountIn,
		AmountOutMin: amountOut,
		Path: []accounts.Address{
			accounts.HexToAddress(PepeContractAddr),
			accounts.HexToAddress(WETH9ContractAddress),
		},
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

	st := mockTrade.BinarySearch(mockPairResp)
	fmt.Println("Max profit:", st.SandwichPrediction.ExpectedProfit)
	fmt.Println("Token sell amount for max profit:", st.SandwichPrediction.SellAmount)
}

func (s *Web3ClientTestSuite) TestSandwichAttackBinSearchV2() {
	amountIn, _ := new(big.Int).SetString("10000000000000000000", 10)
	amountOut, _ := new(big.Int).SetString("235745150537147960000", 10)
	// 1% slippage, meaning they're willing to receive 1% less than the amountOut as minimum acceptable amount
	slippage := new(big.Int).Div(amountOut, big.NewInt(100))
	fmt.Println("slippage", slippage.String())
	amountOutMin := new(big.Int).Add(amountOut, slippage)
	//slippageMargin := new(big.Int).Div(amountOut, big.NewInt(10000))
	//amountOutMinWithMargin := new(big.Int).Add(amountOutMin, slippageMargin)
	fmt.Println("amountOutMin", amountOutMin)
	mockTrade := SwapExactETHForTokensParams{
		Value:        amountIn,
		AmountOutMin: amountOut,
		Path: []accounts.Address{
			accounts.HexToAddress(WETH9ContractAddress),
			accounts.HexToAddress(PepeContractAddr),
		},
	}
	reserve0, _ := new(big.Int).SetString("47956013761392256000", 10)
	reserve1, _ := new(big.Int).SetString("1383382537550055000000", 10)
	token0Addr, token1Addr := StringsToAddresses(WETH9ContractAddress, PepeContractAddr)
	mockPairResp := UniswapV2Pair{
		KLast:    big.NewInt(0),
		Token0:   token0Addr,
		Token1:   token1Addr,
		Reserve0: reserve0,
		Reserve1: reserve1,
	}
	st := mockTrade.BinarySearch(mockPairResp)
	fmt.Println("Max profit:", st.SandwichPrediction.ExpectedProfit)
	fmt.Println("Token sell amount for max profit:", st.SandwichPrediction.SellAmount)
}
func (s *Web3ClientTestSuite) TestSandwichAttackBinSearchV3() {
	amountIn, _ := new(big.Int).SetString("10000000000000000000", 10)
	amountOut, _ := new(big.Int).SetString("235745150537147960000", 10)
	// 1% slippage, meaning they're willing to receive 1% less than the amountOut as minimum acceptable amount
	slippage := new(big.Int).Div(amountOut, big.NewInt(100))
	fmt.Println("slippage", slippage.String())
	amountOutMin := new(big.Int).Add(amountOut, slippage)
	//slippageMargin := new(big.Int).Div(amountOut, big.NewInt(10000))
	//amountOutMinWithMargin := new(big.Int).Add(amountOutMin, slippageMargin)
	fmt.Println("amountOutMin", amountOutMin)
	mockTrade := SwapExactTokensForTokensParams{
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
	st := mockTrade.BinarySearch(mockPairResp)
	fmt.Println("Max profit:", st.SandwichPrediction.ExpectedProfit)
	fmt.Println("Token sell amount for max profit:", st.SandwichPrediction.SellAmount)
}
