package olympus_beacon_cookbooks

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	olympus_cookbooks "github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	zeus_client_oly "github.com/zeus-fyi/olympus/pkg/zeus/client"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/internal_reqs"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type BeaconsTestCookbookTestSuite struct {
	test_suites_base.TestSuite
	ZeusTestClient    zeus_client.ZeusClient
	ZeusTestClientOly zeus_client_oly.ZeusClient
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
	s1 := "spaces-auth"
	s2 := "spaces-key"
	s3 := "age-auth"
	mainnetBeaconCtxNsTop := kns.TopologyKubeCtxNs{
		TopologyID: 0,
		CloudCtxNs: GetBeaconCloudCtxNs(hestia_req_types.Mainnet),
	}

	req := internal_reqs.InternalSecretsCopyFromTo{
		SecretNames: []string{s1, s2, s3},
		FromKns: kns.TopologyKubeCtxNs{
			TopologyID: 0,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "do",
				Region:        "sfo3",
				Context:       "do-sfo3-dev-do-sfo3-zeus",
				Namespace:     "zeus",
				Env:           "dev",
			},
		},
		ToKns: mainnetBeaconCtxNsTop,
	}

	err = t.ZeusTestClientOly.CopySecretsFromToNamespace(context.Background(), req)
	t.Require().Nil(err)
}
func (t *BeaconsTestCookbookTestSuite) TestMainnetClusterDestroy() {
	olympus_cookbooks.ChangeToCookbookDir()
	d := zeus_req_types.TopologyDeployRequest{
		CloudCtxNs: GetBeaconCloudCtxNs(hestia_req_types.Mainnet),
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
	t.ZeusTestClientOly = zeus_client_oly.NewDefaultZeusClient(tc.Bearer)
}

func TestBeaconsTestCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconsTestCookbookTestSuite))
}
