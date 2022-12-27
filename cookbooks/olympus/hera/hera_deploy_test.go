package hera_olympus_cookbook

import (
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

var clusterClassName = "olympus"

var cd = zeus_req_types.ClusterTopologyDeployRequest{
	// ethereumBeacons is a reserved keyword, to make it global, you can replace the below with your own setup
	ClusterClassName:    clusterClassName,
	SkeletonBaseOptions: []string{"hera"},
	CloudCtxNs:          HeraCloudCtxNs,
}

func (t *HeraCookbookTestSuite) TestDeploy() {
	resp, err := t.ZeusTestClient.DeployCluster(ctx, cd)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *HeraCookbookTestSuite) TestChartUpload() {
	resp, err := t.ZeusTestClient.UploadChart(ctx, HeraChartPath, HeraUploadChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)

	HeraDeployKnsReq.TopologyID = resp.TopologyID
	tar := zeus_req_types.TopologyRequest{TopologyID: HeraDeployKnsReq.TopologyID}
	chartResp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)

	err = chartResp.PrintWorkload(HeraChartPath)
	t.Require().Nil(err)
}

func (t *HeraCookbookTestSuite) TestCreateClusterBase() {
	basesInsert := []string{"hera"}
	cc := zeus_req_types.TopologyCreateOrAddComponentBasesToClassesRequest{
		ClusterClassName:   clusterClassName,
		ComponentBaseNames: basesInsert,
	}
	_, err := t.ZeusTestClient.AddComponentBasesToClass(ctx, cc)
	t.Require().Nil(err)
}

func (t *HeraCookbookTestSuite) TestCreateClusterSkeletonBases() {
	cc := zeus_req_types.TopologyCreateOrAddSkeletonBasesToClassesRequest{
		ClusterClassName:  clusterClassName,
		ComponentBaseName: "hera",
		SkeletonBaseNames: []string{"hera"},
	}
	_, err := t.ZeusTestClient.AddSkeletonBasesToClass(ctx, cc)
	t.Require().Nil(err)
}
