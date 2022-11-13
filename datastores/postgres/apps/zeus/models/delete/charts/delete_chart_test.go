package delete_charts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type DeleteChartTestSuite struct {
	test_suites.PGTestSuite
}

func (d *DeleteChartTestSuite) TestDeleteChart() {

	dc := DeleteChart{
		ChartID: 5800802346331576808,
		Deleted: false,
	}
	ctx := context.Background()
	err := dc.DeleteChart(ctx)
	d.Require().Nil(err)
	d.Assert().True(dc.Deleted)
}

func TestDeleteChartTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteChartTestSuite))
}
