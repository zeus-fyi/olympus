package delete_chart_from_topology

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type DeleteChartTestSuite struct {
	test_suites.PGTestSuite
}

func (d *DeleteChartTestSuite) TestDeleteChartFromTopology() {
	dc := DeleteChartFromTopology{
		TopologyID: 1668063792629755904,
		OrgID:      7138983863666903883,
		UserID:     7138958574876245567,
		Deleted:    false,
	}
	ctx := context.Background()
	err := dc.DeleteChartFromTopology(ctx)
	d.Require().Nil(err)
	d.Assert().True(dc.Deleted)
}

func TestDeleteChartTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteChartTestSuite))
}
