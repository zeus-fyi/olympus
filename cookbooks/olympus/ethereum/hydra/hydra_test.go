package olympus_hydra_cookbooks

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	olympus_cookbooks "github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
)

type HydraCookbookTestSuite struct {
	test_suites_base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

var ctx = context.Background()

func (t *HydraCookbookTestSuite) TestClusterDeploy() {
	olympus_cookbooks.ChangeToCookbookDir()

	cdCfg := HydraClusterConfig(&HydraClusterDefinition, "ephemery")
	_, err := cdCfg.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(err)

	cdep := cdCfg.GenerateDeploymentRequest()

	_, err = t.ZeusTestClient.DeployCluster(ctx, cdep)
	t.Require().Nil(err)
}

func (t *HydraCookbookTestSuite) TestClusterDestroy() {
	d := zeus_req_types.TopologyDeployRequest{
		CloudCtxNs: ValidatorCloudCtxNs,
	}
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, d)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *HydraCookbookTestSuite) TestClusterRegisterDefinitions() {
	cdCfg := HydraClusterConfig(&HydraClusterDefinition, "ephemery")
	cd := cdCfg.BuildClusterDefinitions()

	err := cd.CreateClusterClassDefinitions(ctx, t.ZeusTestClient)
	t.Require().Nil(err)
}

func (t *HydraCookbookTestSuite) TestClusterSetup() {
	olympus_cookbooks.ChangeToCookbookDir()

	cd := HydraClusterConfig(&HydraClusterDefinition, "ephemery")
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

func (t *HydraCookbookTestSuite) SetupTest() {
	olympus_cookbooks.ChangeToCookbookDir()

	tc := api_configs.InitLocalTestConfigs()
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
}

func TestHydraCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(HydraCookbookTestSuite))
}
