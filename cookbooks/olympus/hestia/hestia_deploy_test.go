package hestia_olympus_cookbook

import (
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

var clusterClassName = "olympus"

var cd = zeus_req_types.ClusterTopologyDeployRequest{
	ClusterClassName:    clusterClassName,
	SkeletonBaseOptions: []string{"hestia"},
	CloudCtxNs:          HestiaCloudCtxNs,
}

func (t *HestiaCookbookTestSuite) TestDeploy() {
	t.TestChartUpload()
	resp, err := t.ZeusTestClient.DeployCluster(ctx, cd)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *HestiaCookbookTestSuite) TestChartUpload() {
	resp, err := t.ZeusTestClient.UploadChart(ctx, HestiaChartPath, HestiaUploadChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)

	HestiaDeployKnsReq.TopologyID = resp.TopologyID
	tar := zeus_req_types.TopologyRequest{TopologyID: HestiaDeployKnsReq.TopologyID}
	chartResp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)

	err = chartResp.PrintWorkload(HestiaChartPath)
	t.Require().Nil(err)
}

func (t *HestiaCookbookTestSuite) TestCreateClusterBase() {
	basesInsert := []string{"hestia"}
	cc := zeus_req_types.TopologyCreateOrAddComponentBasesToClassesRequest{
		ClusterClassName:   clusterClassName,
		ComponentBaseNames: basesInsert,
	}
	_, err := t.ZeusTestClient.AddComponentBasesToClass(ctx, cc)
	t.Require().Nil(err)
}

func (t *HestiaCookbookTestSuite) TestCreateClusterSkeletonBases() {
	cc := zeus_req_types.TopologyCreateOrAddSkeletonBasesToClassesRequest{
		ClusterClassName:  clusterClassName,
		ComponentBaseName: "hestia",
		SkeletonBaseNames: []string{"hestia"},
	}
	_, err := t.ZeusTestClient.AddSkeletonBasesToClass(ctx, cc)
	t.Require().Nil(err)
}
