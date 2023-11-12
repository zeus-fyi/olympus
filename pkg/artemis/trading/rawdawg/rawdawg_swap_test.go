package artemis_rawdawg_contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

func mockedTrade() *artemis_trading_types.TradeOutcome {
	// TODO, mock some bs trade outcome
	to := &artemis_trading_types.TradeOutcome{
		AmountIn:      artemis_eth_units.Ether,
		AmountInAddr:  artemis_trading_constants.WETH9ContractAddressAccount,
		AmountOut:     artemis_eth_units.NewBigInt(0),
		AmountOutAddr: accounts.HexToAddress("0x8647Ae4E646cd3CE37FdEB4591b0A7928254bb73"),
	}
	return to
}

func (s *ArtemisTradingContractsTestSuite) mockConditions(w3a web3_actions.Web3Actions) common.Address {
	rawDawgAddr := s.testDeployRawdawgContract(w3a)
	to := mockedTrade()
	err := w3a.SetBalanceAtSlotNumber(ctx, to.AmountInAddr.Hex(), rawDawgAddr.Hex(), 3, artemis_eth_units.EtherMultiple(100))
	s.Require().Nil(err)
	rawDawgWethBal, err := w3a.ReadERC20TokenBalance(ctx, to.AmountInAddr.Hex(), rawDawgAddr.Hex())
	s.Require().Nil(err)
	s.Require().Equal(artemis_eth_units.EtherMultiple(100), rawDawgWethBal)
	return rawDawgAddr
}

func (s *ArtemisTradingContractsTestSuite) TestRawDawgExecSimSwap() {
	sessionID := fmt.Sprintf("%s-%s", "mainnet-fork-session", uuid.New().String())
	//sessionID = fmt.Sprintf("%s-%s", "local-network-session", "12b5d9ce-29dd-4f95-8e89-fed4aef2193d")
	w3a := CreateMainnetForkUser(ctx, s.Tc.ProductionLocalTemporalBearerToken, sessionID)
	defer func(sessionID string) {
		err := w3a.EndAnvilSession()
		s.Require().Nil(err)
	}(sessionID)

	s.testRawDawgExecV2Swap(w3a, s.mockConditions(w3a), mockedTrade())
}

// / 0x3B5A9789Fe2c302420b70e5697D8c7015E475b7A
func (s *ArtemisTradingContractsTestSuite) testRawDawgExecV2Swap(w3a web3_actions.Web3Actions, rawDawgAddr common.Address, to *artemis_trading_types.TradeOutcome) {
	scPayload, err := GetRawDawgV2SimSwapAbiPayload(ctx, rawDawgAddr.Hex(), to)
	s.Assert().NotEmpty(scPayload)
	tx, err := w3a.CallFunctionWithData(ctx, scPayload, scPayload.Data)
	s.Assert().Nil(err)
	s.Assert().NotNil(tx)

	owner, err := w3a.GetOwner(ctx, scPayload.ContractABI, rawDawgAddr.Hex())
	s.Require().Nil(err)
	fmt.Println(owner.String())

	//
	//err = w3a.SetBalanceAtSlotNumber(ctx, to.AmountInAddr.String(), rawDawgAddr.String(), 3, artemis_eth_units.EtherMultiple(100))
	//s.Require().Nil(err)
	//
	//rawDawgWethBal, err := w3a.ReadERC20TokenBalance(ctx, to.AmountInAddr.String(), rawDawgAddr.String())
	//s.Require().Nil(err)
	//s.Require().Equal(artemis_eth_units.EtherMultiple(100), rawDawgWethBal)
}
