package artemis_rawdawg_contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

func (s *ArtemisTradingContractsTestSuite) TestRawDawgSimOutUtil() {
	sessionID := fmt.Sprintf("%s-%s", "local-network-session", uuid.New().String())
	//sessionID = fmt.Sprintf("%s-%s", "local-network-session", "12b5d9ce-29dd-4f95-8e89-fed4aef2193d")
	w3a := CreateLocalUser(ctx, s.Tc.ProductionLocalTemporalBearerToken, sessionID)
	defer func(sessionID string) {
		err := w3a.EndAnvilSession()
		s.Require().Nil(err)
	}(sessionID)

	rawdawgAddr := s.testDeployRawdawgContract(w3a)
	s.testRawDawgSimOutUtil(w3a, rawdawgAddr, nil)
}

func (s *ArtemisTradingContractsTestSuite) testRawDawgSimOutUtil(w3a web3_actions.Web3Actions, rawDawgAddr common.Address, to *artemis_trading_types.TradeOutcome) {
	scPayload := GetRawdawgV2SimSwapAbiPayload(rawDawgAddr.String(), to)
	s.Assert().NotEmpty(scPayload)

	tx, err := w3a.CallConstantFunction(ctx, scPayload)
	s.Assert().Nil(err)
	s.Assert().NotNil(tx)
}
