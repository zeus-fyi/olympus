package olympus_hydra_cookbooks

import (
	"fmt"

	olympus_cookbooks "github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
)

func (t *HydraCookbookTestSuite) TestGoerliClusterDeploy() {
	olympus_cookbooks.ChangeToCookbookDir()
	cdCfg := HydraClusterConfig(&HydraClusterDefinition, "goerli")
	//cdCfg.FilterSkeletonBaseUploads = &strings_filter.FilterOpts{
	//	StartsWith: "mev",
	//}
	t.Require().Equal("goerli-staking", cdCfg.CloudCtxNs.Namespace)
	_, err := cdCfg.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(err)

	cdep := cdCfg.GenerateDeploymentRequest()
	_, err = t.ZeusTestClient.DeployCluster(ctx, cdep)
	t.Require().Nil(err)
}

func (t *HydraCookbookTestSuite) TestGoerliClusterDestroy() {
	d := zeus_req_types.TopologyDeployRequest{
		CloudCtxNs: ValidatorCloudCtxNs,
	}
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, d)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *HydraCookbookTestSuite) TestGoerliClusterSetup() {
	olympus_cookbooks.ChangeToCookbookDir()

	cd := HydraClusterConfig(&HydraClusterDefinition, "goerli")
	gcd := cd.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	gdr := cd.GenerateDeploymentRequest()
	t.Assert().NotEmpty(gdr)
	fmt.Println(gdr)

	sbDefs, err := cd.GenerateSkeletonBaseCharts()
	t.Require().Nil(err)
	t.Assert().NotEmpty(sbDefs)
}

func (t *HydraCookbookTestSuite) TestGoerliClusterRegisterDefinitions() {
	cdCfg := HydraClusterConfig(&HydraClusterDefinition, "goerli")
	cd := cdCfg.BuildClusterDefinitions()

	err := cd.CreateClusterClassDefinitions(ctx, t.ZeusTestClient)
	t.Require().Nil(err)
}
