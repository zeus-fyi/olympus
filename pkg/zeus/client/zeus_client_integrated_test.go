package zeus_client

import (
	"time"

	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

func (t *ZeusClientTestSuite) TestChartUploadAndRead() {
	topologyID := t.TestChartUpload()
	deployKnsReq.TopologyID = topologyID
	t.TestReadDemoChart()
}

func (t *ZeusClientTestSuite) TestEndToEnd() {
	topologyID := t.TestChartUpload()
	deployKnsReq.TopologyID = topologyID
	t.TestReadDemoChart()
	resp, err := t.ZeusTestClient.Deploy(ctx, deployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	tar := zeus_req_types.TopologyCloudCtxNsQueryRequest{
		CloudCtxNs: topCloudCtxNs,
	}
	workloadDeployed, err := t.ZeusTestClient.ReadNamespaceWorkload(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(workloadDeployed)

	t.Assert().NotEmpty(workloadDeployed.DeploymentList)
	t.Assert().NotEmpty(workloadDeployed.ServiceList)

	destroyResp, err := t.ZeusTestClient.DestroyDeploy(ctx, deployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(destroyResp)

	time.Sleep(10 * time.Second)
	workloadPostDestroy, err := t.ZeusTestClient.ReadNamespaceWorkload(ctx, tar)
	t.Require().Nil(err)
	t.Assert().Len(workloadPostDestroy.DeploymentList.Items, 0)
	t.Assert().Len(workloadPostDestroy.ServiceList.Items, 0)
}
