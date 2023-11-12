package artemis_rawdawg_contract

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

func (s *ArtemisTradingContractsTestSuite) TestRawDawgSimOutUtil() {
	sessionID := fmt.Sprintf("%s-%s", "forked-mainnet-session", uuid.New().String())
	w3a := CreateUser(ctx, "mainnet", s.Tc.ProductionLocalTemporalBearerToken, sessionID)
	defer func(sessionID string) {
		err := w3a.EndAnvilSession()
		s.Require().Nil(err)
	}(sessionID)

	rdAddr, abiFile := s.mockConditions(w3a, mockedTrade())
	s.testRawDawgExecV2SwapSimMainnet(w3a, rdAddr, abiFile, mockedTrade())
}

func (s *ArtemisTradingContractsTestSuite) testRawDawgExecV2SwapSimMainnet(w3a web3_actions.Web3Actions, rawDawgAddr common.Address, abiFile *abi.ABI, to *artemis_trading_types.TradeOutcome) {
	scPayload := GetRawDawgV2SimSwapAbiPayload(ctx, rawDawgAddr.Hex(), abiFile, to)
	s.Assert().NotEmpty(scPayload)
	resp, err := w3a.CallConstantFunction(ctx, scPayload)
	s.Assert().Nil(err)
	s.Assert().NotNil(resp)

	for _, val := range resp {
		fmt.Println(val)

		bgn, ok := val.(big.Int)
		if ok {
			fmt.Println(bgn.String())
		}
	}
}

/*
   // Calculate the gas used
      gasUsed = gasBefore - gasleft();

      // Calculate the simulated final balances
      balanceTokenInAfter = IERC20(_token_in).balanceOf(address(this));
      balanceTokenOutAfter = IERC20(_token_out).balanceOf(address(this));
*/
