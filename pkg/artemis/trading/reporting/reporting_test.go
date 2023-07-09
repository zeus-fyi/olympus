package artemis_reporting

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
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
	rw, err := GetRewardsHistory(ctx)
	s.Assert().Nil(err)
	s.Assert().NotNil(rw)

	total := artemis_eth_units.NewBigInt(0)
	totalWithoutNegatives := artemis_eth_units.NewBigInt(0)

	for _, v := range rw.Map {
		//fmt.Println(k)
		total = artemis_eth_units.AddBigInt(total, v.ExpectedProfitAmountOut)
		fmt.Println(v.ExpectedProfitAmountOut)
		if artemis_eth_units.IsXGreaterThanY(v.ExpectedProfitAmountOut, artemis_eth_units.NewBigInt(0)) {
			totalWithoutNegatives = artemis_eth_units.AddBigInt(totalWithoutNegatives, v.ExpectedProfitAmountOut)
		}
	}

	fmt.Println("total eth profit", total.String())
	fmt.Println("total eth profit without negatives", totalWithoutNegatives.String())
}

func (s *ReportingTestSuite) Test1() {

}

func TestReportingTestSuite(t *testing.T) {
	suite.Run(t, new(ReportingTestSuite))
}
