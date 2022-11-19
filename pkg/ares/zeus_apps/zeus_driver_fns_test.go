package zeus_apps

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

func (t *ZeusAppsTestSuite) TestReadDemoChart(topID int, chartPath filepaths.Path) {
	tar := zeus_req_types.TopologyRequest{TopologyID: topID}
	resp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	// prints the chart output for inspection
	err = resp.PrintWorkload(chartPath)
	t.Require().Nil(err)
}
