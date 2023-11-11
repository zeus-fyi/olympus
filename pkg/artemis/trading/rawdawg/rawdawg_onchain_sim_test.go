package artemis_rawdawg_contract

import artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"

func (s *ArtemisTradingContractsTestSuite) TestRawDawgSimOutUtil() {
	tradingSwapContractAddr := ""
	pairContractAddr := ""
	to := &artemis_trading_types.TradeOutcome{}
	tmp := GetRawdawgSwapAbiPayload(tradingSwapContractAddr, pairContractAddr, to, false)
	s.Assert().NotEmpty(tmp)
}
