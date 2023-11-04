package artemis_reporting

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

var ctx = context.Background()

type ReportingTestSuite struct {
	w3c        web3_client.Web3Client
	w3cArchive web3_client.Web3Client
	test_suites_encryption.EncryptionTestSuite
}

func (s *ReportingTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.w3c = web3_client.NewWeb3ClientFakeSigner("https://eth.zeus.fyi")
	s.w3c.AddBearerToken(s.Tc.ProductionLocalBearerToken)
	s.w3cArchive = web3_client.NewWeb3ClientFakeSigner("https://eth-mainnet.g.alchemy.com/v2/cdVqiD1oZGvBiNEU8rDYt5kb6Q24nBMB")
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}

func (s *ReportingTestSuite) TestCalculateProfits() {
	// init 17639300, current 17658962
	//  artemis_trading_constants.V2SwapExactIn
	rhf := RewardHistoryFilter{
		FromBlock:   17658962,
		TradeMethod: "any",
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
	var addresses []string
	for _, v := range historySlice {
		total1 := v.Count + v.FailedCount
		if v.Count < 2 || (float64(v.FailedCount) > float64(total1)*0.2) {
			// v.FailedCount is more than 10% of the total
			log.Info().Str("token", v.AmountOutToken.Name()).Str("address", v.AmountOutToken.Address.String()).Int("successCount", v.Count).Int("failureCount", v.FailedCount).Msg("failed count is more than 20% of the total")
			continue
		}

		if artemis_eth_units.IsXGreaterThanY(artemis_eth_units.NewBigInt(0), rw.Map[v.AmountOutToken.Address.String()].ExpectedProfitAmountOut) {
			log.Warn().Str("token", v.AmountOutToken.Name()).Str("address", v.AmountOutToken.Address.String()).Int("successCount", v.Count).Int("failureCount", v.FailedCount).Msg("expected profit is negative")
			negCount++
			continue
		}
		fmt.Println(
			"successCount", v.Count, "failedCount", v.FailedCount, "expProfits", v.ExpectedProfitAmountOut.String(),
			v.AmountOutToken.Name(), v.AmountOutToken.Address.String(),
			"num", v.AmountOutToken.TransferTax.Numerator.String(), "den", v.AmountOutToken.TransferTax.Denominator.String())
		addresses = append(addresses, v.AmountOutToken.Address.String())
	}

	fmt.Println("total eth profit", artemis_eth_units.DivBigIntToFloat(total, artemis_eth_units.Ether).String())
	fmt.Println("negatives", negCount)
	fmt.Println("total eth profit without negatives", artemis_eth_units.DivBigIntToFloat(totalWithoutNegatives, artemis_eth_units.Ether).String())
	//
	//err = artemis_risk_analysis.SetTradingPermission(ctx, addresses, 1, true)
	//s.Assert().Nil(err)
}

func (s *ReportingTestSuite) TestListSimResults() {
	// init 17639300, current 17658962
	//  artemis_trading_constants.V2SwapExactIn
	rhf := RewardHistoryFilter{
		FromBlock:   17658962,
		TradeMethod: "any",
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

	for _, v := range historySlice {
		fmt.Println(
			"successCount", v.Count, "failedCount", v.FailedCount, "expProfits", v.ExpectedProfitAmountOut.String(),
			v.AmountOutToken.Name(), v.AmountOutToken.Address.String(),
			"num", v.AmountOutToken.TransferTax.Numerator.String(), "den", v.AmountOutToken.TransferTax.Denominator.String())
	}
}

func TestReportingTestSuite(t *testing.T) {
	suite.Run(t, new(ReportingTestSuite))
}
