package artemis_reporting

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

var ctx = context.Background()

type ReportingTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
}

func (s *ReportingTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}

func (s *ReportingTestSuite) TestCalculateProfits() {
	// init 17639300, current 17658962
	//  artemis_trading_constants.V2SwapExactIn
	rhf := RewardHistoryFilter{
		FromBlock:   17658962,
		TradeMethod: artemis_trading_constants.V2SwapExactIn,
	}
	rw, err := GetRewardsHistory(ctx, rhf)
	s.Assert().Nil(err)
	s.Assert().NotNil(rw)

	total := artemis_eth_units.NewBigInt(0)
	totalWithoutNegatives := artemis_eth_units.NewBigInt(0)
	// Sort slice
	var historySlice []RewardsHistory
	for _, v := range rw.Map {
		total = artemis_eth_units.AddBigInt(total, v.ExpectedProfitAmountOut)
		if artemis_eth_units.IsXGreaterThanY(v.ExpectedProfitAmountOut, artemis_eth_units.NewBigInt(0)) {
			totalWithoutNegatives = artemis_eth_units.AddBigInt(totalWithoutNegatives, v.ExpectedProfitAmountOut)
		}
		rh := RewardsHistory{
			FailedCount:             v.FailedCount,
			AmountOutToken:          v.AmountOutToken,
			Count:                   v.Count,
			ExpectedProfitAmountOut: v.ExpectedProfitAmountOut,
		}
		historySlice = append(historySlice, rh)
	}

	sort.SliceStable(historySlice, func(i, j int) bool {
		// Descending order
		return historySlice[i].Count > historySlice[j].Count
	})

	negCount := 0
	for _, v := range historySlice {
		if v.Count < 1 {
			continue
		}
		if artemis_eth_units.IsXGreaterThanY(artemis_eth_units.NewBigInt(0), rw.Map[v.AmountOutToken.Address.String()].ExpectedProfitAmountOut) {
			negCount++
			continue
		}
		fmt.Println(
			"successCount", v.Count, "failedCount", v.FailedCount, "expProfits", v.ExpectedProfitAmountOut.String(),
			v.AmountOutToken.Name(), v.AmountOutToken.Address.String(),
			"num", v.AmountOutToken.TransferTax.Numerator.String(), "den", v.AmountOutToken.TransferTax.Denominator.String())
	}

	fmt.Println("total eth profit", artemis_eth_units.DivBigIntToFloat(total, artemis_eth_units.Ether).String())
	fmt.Println("negatives", negCount)
	fmt.Println("total eth profit without negatives", artemis_eth_units.DivBigIntToFloat(totalWithoutNegatives, artemis_eth_units.Ether).String())
}

func (s *ReportingTestSuite) Test1() {

}

func TestReportingTestSuite(t *testing.T) {
	suite.Run(t, new(ReportingTestSuite))
}
