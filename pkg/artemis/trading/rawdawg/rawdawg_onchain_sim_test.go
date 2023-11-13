package artemis_rawdawg_contract

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

func (s *ArtemisTradingContractsTestSuite) TestRawDawgSimOutUtil() {
	sessionOne := fmt.Sprintf("%s-%s", "forked-mainnet-session-1", uuid.New().String())
	w3a := CreateUser(ctx, "mainnet", s.Tc.ProductionLocalTemporalBearerToken, sessionOne)
	s.T().Cleanup(func() {
		func(sessionID string) {
			fmt.Printf("CLEANUP: ENDING SESSION %s", sessionID)
			err := w3a.EndAnvilSession()
			s.Require().Nil(err)
		}(sessionOne)
	})
	rdAddr, abiFile := s.mockConditions(w3a, mockedTrade())
	s.testRawDawgExecV2SwapSimMainnet(w3a, rdAddr, abiFile, mockedTrade(), true)

	/*
		99957705576280962606
		957705576280962606
	*/

	/*
		55925319574428105816755167200
		58049694401678927436191764679
	*/

	sessionTwo := fmt.Sprintf("%s-%s", "forked-mainnet-session-2", uuid.New().String())
	w3a2 := CreateUser(ctx, "mainnet", s.Tc.ProductionLocalTemporalBearerToken, sessionTwo)
	s.T().Cleanup(func() {
		func(sessionID string) {
			fmt.Printf("CLEANUP: ENDING SESSION %s", sessionID)
			err := w3a.EndAnvilSession()
			s.Require().Nil(err)
		}(sessionTwo)
	})

	//s.testRawDawgExecV2SwapSimMainnet(w3a2, rdAddr, abiFile, mockedTrade(), true)
	s.testRawDawgExecV2SwapMainnet(w3a2, rdAddr, abiFile, mockedTrade())
}

func (s *ArtemisTradingContractsTestSuite) testRawDawgExecV2SwapSimMainnet(w3a web3_actions.Web3Actions, rawDawgAddr common.Address, abiFile *abi.ABI, to *artemis_trading_types.TradeOutcome, buyAndSell bool) {
	fmt.Println("SIM")
	ao := artemis_eth_units.NewBigIntFromStr("55925319574428105816755167200")

	to = &artemis_trading_types.TradeOutcome{
		AmountIn:      artemis_eth_units.EtherMultiple(1),
		AmountInAddr:  artemis_trading_constants.WETH9ContractAddressAccount,
		AmountOut:     ao,
		AmountOutAddr: artemis_trading_constants.BoboTokenAddressAccount,
	}
	var scPayload *web3_actions.SendContractTxPayload
	if buyAndSell {
		scPayload = GetRawDawgV2SimSwapBuySellAbiPayload(ctx, rawDawgAddr.Hex(), abiFile, to)
	} else {
		scPayload = GetRawDawgV2SimSwapAbiPayload(ctx, rawDawgAddr.Hex(), abiFile, to)
	}
	s.Assert().NotEmpty(scPayload)
	resp, err := w3a.CallConstantFunction(ctx, scPayload)
	s.Assert().Nil(err)
	s.Assert().NotNil(resp)

	rawDawgTokenBal, err := w3a.ReadERC20TokenBalance(ctx, to.AmountInAddr.Hex(), rawDawgAddr.Hex())
	s.Require().Nil(err)
	fmt.Println("tokenOut", rawDawgTokenBal.String())

	ethBa, err := w3a.GetBalance(ctx, rawDawgAddr.Hex(), nil)
	s.Require().Nil(err)
	fmt.Println("ethBal", ethBa.String())
	//
	for _, val := range resp {
		fmt.Println(val)
		bgn, ok := val.(big.Int)
		if ok {
			fmt.Println(bgn.String())
		}
	}

	rawDawgTokenBal, err = w3a.ReadERC20TokenBalance(ctx, to.AmountOutAddr.Hex(), rawDawgAddr.Hex())
	s.Require().Nil(err)
	fmt.Println(rawDawgTokenBal.String())
}

func (s *ArtemisTradingContractsTestSuite) testRawDawgExecV2SwapMainnet(w3a web3_actions.Web3Actions, rawDawgAddr common.Address, abiFile *abi.ABI, to *artemis_trading_types.TradeOutcome) {
	fmt.Println("SWAP")
	rawDawgTokenBal, err := w3a.ReadERC20TokenBalance(ctx, to.AmountInAddr.Hex(), rawDawgAddr.Hex())
	s.Require().Nil(err)
	fmt.Println(rawDawgTokenBal.String())

	ethBa, err := w3a.GetBalance(ctx, rawDawgAddr.Hex(), nil)
	s.Require().Nil(err)
	fmt.Println(ethBa.String())

	ao := artemis_eth_units.NewBigIntFromStr("55925319574428105816755167200")

	to = &artemis_trading_types.TradeOutcome{
		AmountIn:      artemis_eth_units.EtherMultiple(1),
		AmountInAddr:  artemis_trading_constants.WETH9ContractAddressAccount,
		AmountOut:     ao,
		AmountOutAddr: artemis_trading_constants.BoboTokenAddressAccount,
	}
	/*
	   113213872122962575389353695009
	   56819205366795356299650473948
	   30702843238000000000000000000
	*/
	fmt.Println("SWAP TOKEN OUT BEFORE")
	rawDawgTokenBal, err = w3a.ReadERC20TokenBalance(ctx, to.AmountOutAddr.Hex(), rawDawgAddr.Hex())
	s.Require().Nil(err)
	fmt.Println(rawDawgTokenBal.String())

	tx, err := ExecSmartContractTradingSwap(ctx, w3a, rawDawgAddr.Hex(), abiFile, to)
	s.Require().Nil(err)
	s.Assert().NotNil(tx)
	fmt.Println("SWAP TOKEN OUT AFTER")

	rawDawgTokenBal, err = w3a.ReadERC20TokenBalance(ctx, to.AmountOutAddr.Hex(), rawDawgAddr.Hex())
	s.Require().Nil(err)
	fmt.Println(rawDawgTokenBal.String())
}

/*
SWAP
100000000000000000000
10000000000000000000000000
0
55925319574428105816755167200
SIM
99000000000000000000
10000000000000000000000000
167237629011979372435346874319
56365952548769869184628169138
55925319574428105816755167200

114162747207475440069759043250
55925319574428105816755167200
57818780273898964994726729840
*/
