package artemis_rawdawg_contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_utils "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/utils"
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

func (s *ArtemisTradingContractsTestSuite) TestMock() {
	to := mockedTrade()
	pairContractAddr, token0, token1 := artemis_utils.CreateV2TradingPair(to.AmountInAddr, to.AmountOutAddr)
	s.Require().NotEmpty(pairContractAddr)
	fmt.Println(pairContractAddr, " ", token0.String(), " ", token1.String())
}

func (s *ArtemisTradingContractsTestSuite) TestRawDawgExecSimSwap() {
	sessionID := fmt.Sprintf("%s-%s", "mainnet-fork-session", uuid.New().String())
	//sessionID = fmt.Sprintf("%s-%s", "local-network-session", "12b5d9ce-29dd-4f95-8e89-fed4aef2193d")

	w3a := CreateLocalUser(ctx, s.Tc.ProductionLocalTemporalBearerToken, sessionID)
	w3a.AddDefaultEthereumMainnetTableHeader()
	defer func(sessionID string) {
		err := w3a.EndAnvilSession()
		s.Require().Nil(err)
	}(sessionID)

	rawDawgAddr := s.testDeployRawdawgContract(w3a)
	to := mockedTrade()
	s.testRawDawgExecV2Swap(w3a, rawDawgAddr, to)
}

// / 0x3B5A9789Fe2c302420b70e5697D8c7015E475b7A
func (s *ArtemisTradingContractsTestSuite) testRawDawgExecV2Swap(w3a web3_actions.Web3Actions, rawDawgAddr common.Address, to *artemis_trading_types.TradeOutcome) {
	scPayload := GetRawdawgV2SimSwapAbiPayload(rawDawgAddr.String(), to)
	s.Assert().NotEmpty(scPayload)
	tx, err := w3a.CallFunctionWithArgs(ctx, scPayload)
	s.Assert().Nil(err)
	s.Assert().NotNil(tx)
}
