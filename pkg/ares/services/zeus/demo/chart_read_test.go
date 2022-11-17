package zeus_demo

import (
	"github.com/zeus-fyi/olympus/pkg/ares/demo"
	read_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/read"
)

func (t *AresDemoTestSuite) TestReadObolDemoChart() {
	tar := read_infra.TopologyReadRequest{TopologyID: obolDeployKnsReq.TopologyID}
	resp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
	demo.ChangeDirToAresDemoDir()

	p := demo.DemoReadChartThenWritePath()
	err = resp.PrintWorkload(p)
	t.Require().Nil(err)
}
