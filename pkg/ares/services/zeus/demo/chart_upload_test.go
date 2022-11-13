package zeus_demo

import (
	"github.com/zeus-fyi/olympus/pkg/ares/demo"
	create_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"
)

func (t *AresDemoTestSuite) TestCreateAndUploadDemoChart() create_infra.TopologyCreateResponse {
	demo.ChangeDirToAresDemoDir()
	p := demo.DemoPath()
	chartInfo := demo.DemoChartUploadRequest()
	resp, err := t.ZeusTestClient.UploadChart(ctx, p, chartInfo)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.ID)
	return resp
}
