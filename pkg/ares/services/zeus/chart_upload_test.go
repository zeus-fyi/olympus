package ares_zeus_driver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/ares/ethereum"
	ares_test_suite "github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/ares"
)

type AresZeusTestSuite struct {
	ares_test_suite.AresTestSuite
}

var ctx = context.Background()

func (t *AresZeusTestSuite) TestCreateAndUploadBeaconChart() {
	ethereum.ChangeDirToAresEthereumDir()
	p := ethereum.ConsensusClientPath()
	chartInfo := ethereum.ConsensusClientChartUploadRequest()
	resp, err := t.ProdZeusClient.UploadChart(ctx, p, chartInfo)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.ID)
}

func TestAresZeusTestSuite(t *testing.T) {
	suite.Run(t, new(AresZeusTestSuite))
}
