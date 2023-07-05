package async_analysis

import (
	"fmt"

	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

// TODO, setup global test suite

func (s *ArtemisRealTimeTradingTestSuite) testEthSimTransferFeeAnalysisSetup() {
	uni := web3_client.InitUniswapClient(ctx, s.UserA)
	doge2 := artemis_trading_constants.Doge2ContractAddr
	s.ca = NewERC20ContractAnalysis(&uni, doge2)
	s.ca.UserB = s.UserB
	err := s.ca.UserA.HardHatResetNetwork(ctx, s.Tc.QuiknodeLiveNode, 17601900)
	s.Require().Nil(err)
}

func (s *ArtemisRealTimeTradingTestSuite) TestCalculateTransferFeeTaxRange() {
	s.testEthSimTransferFeeAnalysisSetup()
	percent, err := s.ca.CalculateTransferFeeTax(ctx, artemis_eth_units.EtherMultiple(1))
	s.Assert().Nil(err)
	fmt.Println(percent.Numerator.String())
	fmt.Println(percent.Denominator.String())

	percent, err = s.ca.CalculateTransferFeeTax(ctx, artemis_eth_units.GweiMultiple(1))
	s.Assert().Nil(err)
	fmt.Println(percent.Numerator.String())
	fmt.Println(percent.Denominator.String())

	percent, err = s.ca.CalculateTransferFeeTax(ctx, artemis_eth_units.EtherMultiple(100))
	s.Assert().Nil(err)
	fmt.Println(percent.Numerator.String())
	fmt.Println(percent.Denominator.String())
}
