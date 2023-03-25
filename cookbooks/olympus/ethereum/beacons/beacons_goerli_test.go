package olympus_beacon_cookbooks

import (
	"fmt"

	olympus_cookbooks "github.com/zeus-fyi/olympus/cookbooks"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
)

func (t *BeaconsTestCookbookTestSuite) TestGoerliClusterDeploy() {
	olympus_cookbooks.ChangeToCookbookDir()
	cdCfg := GoerliBeaconBaseClusterDefinition
	//cdCfg.FilterSkeletonBaseUploads = &strings_filter.FilterOpts{
	//	StartsWith: "ingress",
	//}
	t.Require().Equal("athena-beacon-goerli", cdCfg.CloudCtxNs.Namespace)
	_, err := cdCfg.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(err)

	cdep := cdCfg.GenerateDeploymentRequest()
	_, err = t.ZeusTestClient.DeployCluster(ctx, cdep)
	t.Require().Nil(err)
}

func (t *BeaconsTestCookbookTestSuite) TestGoerliClusterDestroy() {
	d := zeus_req_types.TopologyDeployRequest{
		CloudCtxNs: GetBeaconCloudCtxNs(hestia_req_types.Goerli),
	}
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, d)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *BeaconsTestCookbookTestSuite) TestGoerliClusterSetup() {
	olympus_cookbooks.ChangeToCookbookDir()

	gcd := GoerliBeaconBaseClusterDefinition.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	gdr := GoerliBeaconBaseClusterDefinition.GenerateDeploymentRequest()
	t.Assert().NotEmpty(gdr)
	fmt.Println(gdr)

	sbDefs, err := GoerliBeaconBaseClusterDefinition.GenerateSkeletonBaseCharts()
	t.Require().Nil(err)
	t.Assert().NotEmpty(sbDefs)
}

func (t *BeaconsTestCookbookTestSuite) TestCreateClusterBase() {
	basesInsert := []string{"beaconIngress"}
	cc := zeus_req_types.TopologyCreateOrAddComponentBasesToClassesRequest{
		ClusterClassName:   GoerliBeaconBaseClusterDefinition.ClusterClassName,
		ComponentBaseNames: basesInsert,
	}
	_, err := t.ZeusTestClient.AddComponentBasesToClass(ctx, cc)
	t.Require().Nil(err)
}

func (t *BeaconsTestCookbookTestSuite) TestCreateClusterSkeletonBases() {
	cc := zeus_req_types.TopologyCreateOrAddSkeletonBasesToClassesRequest{
		ClusterClassName:  GoerliBeaconBaseClusterDefinition.ClusterClassName,
		ComponentBaseName: "beaconIngress",
		SkeletonBaseNames: []string{"ingress"},
	}
	_, err := t.ZeusTestClient.AddSkeletonBasesToClass(ctx, cc)
	t.Require().Nil(err)
}

func (t *BeaconsTestCookbookTestSuite) TestGoerliClusterRegisterDefinitions() {
	cd := GoerliBeaconBaseClusterDefinition.BuildClusterDefinitions()
	err := cd.CreateClusterClassDefinitions(ctx, t.ZeusTestClient)
	t.Require().Nil(err)
}
