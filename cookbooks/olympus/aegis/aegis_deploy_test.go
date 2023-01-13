package aegis_olympus_cookbook

import "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"

const (
	aegis   = "aegis"
	olympus = "olympus"
)

var cd = zeus_req_types.ClusterTopologyDeployRequest{
	ClusterClassName:    olympus,
	SkeletonBaseOptions: []string{aegis},
	CloudCtxNs:          AegisCloudCtxNs,
}

func (t *AegisCookbookTestSuite) TestDeploy() {
	resp, err := t.ZeusTestClient.DeployCluster(ctx, cd)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	//t.TestAegisSecretsCopy()
}

func (t *AegisCookbookTestSuite) TestChartUpload() {
	resp, err := t.ZeusTestClient.UploadChart(ctx, AegisChartPath, AegisUploadChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)

	AegisDeployKnsReq.TopologyID = resp.TopologyID
	tar := zeus_req_types.TopologyRequest{TopologyID: AegisDeployKnsReq.TopologyID}
	chartResp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)

	err = chartResp.PrintWorkload(AegisChartPath)
	t.Require().Nil(err)
}

func (t *AegisCookbookTestSuite) TestCreateClusterBase() {
	basesInsert := []string{aegis}
	cc := zeus_req_types.TopologyCreateOrAddComponentBasesToClassesRequest{
		ClusterClassName:   olympus,
		ComponentBaseNames: basesInsert,
	}
	_, err := t.ZeusTestClient.AddComponentBasesToClass(ctx, cc)
	t.Require().Nil(err)
}

func (t *AegisCookbookTestSuite) TestCreateClusterSkeletonBases() {
	cc := zeus_req_types.TopologyCreateOrAddSkeletonBasesToClassesRequest{
		ClusterClassName:  olympus,
		ComponentBaseName: aegis,
		SkeletonBaseNames: []string{aegis},
	}
	_, err := t.ZeusTestClient.AddSkeletonBasesToClass(ctx, cc)
	t.Require().Nil(err)
}
