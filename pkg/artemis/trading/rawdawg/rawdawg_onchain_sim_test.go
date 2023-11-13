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
	sessionID := fmt.Sprintf("%s-%s", "forked-mainnet-session", uuid.New().String())
	fmt.Println(sessionID)
	w3a := CreateUser(ctx, "mainnet", s.Tc.ProductionLocalTemporalBearerToken, sessionID)
	s.T().Cleanup(func() {
		func(sessionID string) {
			err := w3a.EndAnvilSession()
			s.Require().Nil(err)
		}(sessionID)
	})
	rdAddr, abiFile := s.mockConditions(w3a, mockedTrade())
	//s.testRawDawgExecV2SwapMainnet(w3a, rdAddr, abiFile, mockedTrade())
	s.testRawDawgExecV2SwapSimMainnet(w3a, rdAddr, abiFile, mockedTrade())
}

func (s *ArtemisTradingContractsTestSuite) testRawDawgExecV2SwapSimMainnet(w3a web3_actions.Web3Actions, rawDawgAddr common.Address, abiFile *abi.ABI, to *artemis_trading_types.TradeOutcome) {
	fmt.Println("SIM")

	to = &artemis_trading_types.TradeOutcome{
		AmountIn:      artemis_eth_units.EtherMultiple(1),
		AmountInAddr:  artemis_trading_constants.WETH9ContractAddressAccount,
		AmountOut:     artemis_eth_units.EtherMultiple(65_702_843_238),
		AmountOutAddr: artemis_trading_constants.BoboTokenAddressAccount,
	}

	scPayload := GetRawDawgV2SimSwapAbiPayload(ctx, rawDawgAddr.Hex(), abiFile, to)
	s.Assert().NotEmpty(scPayload)
	resp, err := w3a.CallConstantFunction(ctx, scPayload)
	s.Assert().Nil(err)
	s.Assert().NotNil(resp)

	rawDawgTokenBal, err := w3a.ReadERC20TokenBalance(ctx, to.AmountInAddr.Hex(), rawDawgAddr.Hex())
	s.Require().Nil(err)
	fmt.Println(rawDawgTokenBal.String())

	ethBa, err := w3a.GetBalance(ctx, rawDawgAddr.Hex(), nil)
	s.Require().Nil(err)
	fmt.Println(ethBa.String())
	//
	for _, val := range resp {
		fmt.Println(val)
		bgn, ok := val.(big.Int)
		if ok {
			fmt.Println(bgn.String())
		}
	}

	/*
		255866077281820760280
		128131335018675421342
	*/
	// 127511844451648205958
	//tx, err := ExecSmartContractTradingSwap(ctx, w3a, rawDawgAddr.Hex(), abiFile, to)
	//s.Require().Nil(err)
	//s.Assert().NotNil(tx)
	//
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
		AmountOutAddr: artemis_trading_constants.MogTokenAddressAccount,
	}
	/*
	   113213872122962575389353695009
	   56819205366795356299650473948
	   30702843238000000000000000000
	*/
	rawDawgTokenBal, err = w3a.ReadERC20TokenBalance(ctx, to.AmountOutAddr.Hex(), rawDawgAddr.Hex())
	s.Require().Nil(err)
	fmt.Println(rawDawgTokenBal.String())

	tx, err := ExecSmartContractTradingSwap(ctx, w3a, rawDawgAddr.Hex(), abiFile, to)
	s.Require().Nil(err)
	s.Assert().NotNil(tx)

	rawDawgTokenBal, err = w3a.ReadERC20TokenBalance(ctx, to.AmountOutAddr.Hex(), rawDawgAddr.Hex())
	s.Require().Nil(err)
	fmt.Println(rawDawgTokenBal.String())
}
