package olympus_beacon_cookbooks

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

type BeaconsTestCookbookTestSuite struct {
	test_suites_base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

var ctx = context.Background()

func (t *BeaconsTestCookbookTestSuite) TestMainnetClusterDeploy() {
	olympus_cookbooks.ChangeToCookbookDir()
	ClusterConfigEnvVars(&MainnetBeaconBaseClusterDefinition, "mainnet")

	_, err := MainnetBeaconBaseClusterDefinition.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(err)

	cdep := MainnetBeaconBaseClusterDefinition.GenerateDeploymentRequest()
	_, err = t.ZeusTestClient.DeployCluster(ctx, cdep)
	t.Require().Nil(err)
}
func (t *BeaconsTestCookbookTestSuite) TestMainnetClusterDestroy() {
	olympus_cookbooks.ChangeToCookbookDir()
	d := zeus_req_types.TopologyDeployRequest{
		CloudCtxNs: MainnetAthenaBeaconCloudCtxNs,
	}
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, d)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *BeaconsTestCookbookTestSuite) TestMainnetClusterSetup() {
	olympus_cookbooks.ChangeToCookbookDir()
	ClusterConfigEnvVars(&MainnetBeaconBaseClusterDefinition, "mainnet")

	gcd := MainnetBeaconBaseClusterDefinition.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	gdr := MainnetBeaconBaseClusterDefinition.GenerateDeploymentRequest()
	t.Assert().NotEmpty(gdr)
	fmt.Println(gdr)

	sbDefs, err := MainnetBeaconBaseClusterDefinition.GenerateSkeletonBaseCharts()
	t.Require().Nil(err)
	t.Assert().NotEmpty(sbDefs)
}
func (t *BeaconsTestCookbookTestSuite) TestMainnetClusterCreate() {
	olympus_cookbooks.ChangeToCookbookDir()
	cdCfg := MainnetBeaconBaseClusterDefinition

	gcd := cdCfg.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)
	err := gcd.CreateClusterClassDefinitions(ctx, t.ZeusTestClient)
	t.Require().Nil(err)

	gdr := cdCfg.GenerateDeploymentRequest()
	t.Assert().NotEmpty(gdr)
	fmt.Println(gdr)

	sbDefs, err := cdCfg.GenerateSkeletonBaseCharts()
	t.Require().Nil(err)
	t.Assert().NotEmpty(sbDefs)
}

func (t *BeaconsTestCookbookTestSuite) TestClusterSetup() {
	olympus_cookbooks.ChangeToCookbookDir()

	gcd := EphemeralBeaconBaseClusterDefinition.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	gdr := EphemeralBeaconBaseClusterDefinition.GenerateDeploymentRequest()
	t.Assert().NotEmpty(gdr)
	fmt.Println(gdr)

	sbDefs, err := EphemeralBeaconBaseClusterDefinition.GenerateSkeletonBaseCharts()
	t.Require().Nil(err)
	t.Assert().NotEmpty(sbDefs)
}

func (t *BeaconsTestCookbookTestSuite) SetupTest() {
	olympus_cookbooks.ChangeToCookbookDir()

	tc := api_configs.InitLocalTestConfigs()
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
}

func TestBeaconsTestCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconsTestCookbookTestSuite))
}
