package hestia_olympus_cookbook

import "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"

var zeusCloudClusterDef = zeus_req_types.ClusterTopologyDeployRequest{
	ClusterClassName:    clusterClassName,
	SkeletonBaseOptions: []string{"zeusCloud"},
	CloudCtxNs:          ZeusCloudCloudCtxNs,
}

// TODO needs to add env vars for api urls
func (t *HestiaCookbookTestSuite) TestDeployZeusCloud() {
	t.TestChartUpload()
	resp, err := t.ZeusTestClient.DeployCluster(ctx, zeusCloudClusterDef)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *HestiaCookbookTestSuite) TestChartUploadZeusCloud() {
	resp, err := t.ZeusTestClient.UploadChart(ctx, ZeusCloudChartPath, ZeusCloudUploadChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)

	ZeusCloudDeployKnsReq.TopologyID = resp.TopologyID
	tar := zeus_req_types.TopologyRequest{TopologyID: ZeusCloudDeployKnsReq.TopologyID}
	chartResp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)

	err = chartResp.PrintWorkload(ZeusCloudChartPath)
	t.Require().Nil(err)
}

func (t *HestiaCookbookTestSuite) TestCreateClusterBaseZeusCloud() {
	basesInsert := []string{"zeusCloud"}
	cc := zeus_req_types.TopologyCreateOrAddComponentBasesToClassesRequest{
		ClusterClassName:   clusterClassName,
		ComponentBaseNames: basesInsert,
	}
	_, err := t.ZeusTestClient.AddComponentBasesToClass(ctx, cc)
	t.Require().Nil(err)
}

func (t *HestiaCookbookTestSuite) TestCreateClusterSkeletonBasesZeusCloud() {
	cc := zeus_req_types.TopologyCreateOrAddSkeletonBasesToClassesRequest{
		ClusterClassName:  clusterClassName,
		ComponentBaseName: "zeusCloud",
		SkeletonBaseNames: []string{"zeusCloud"},
	}
	_, err := t.ZeusTestClient.AddSkeletonBasesToClass(ctx, cc)
	t.Require().Nil(err)
}
