package olympus_hardhat

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	olympus_cookbooks "github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
)

var ctx = context.Background()

func (t *ZeusP2PCrawlerTestSuite) TestClusterDeploy() {
	olympus_cookbooks.ChangeToCookbookDir()
	cdep := hardhatClusterDefinition.GenerateDeploymentRequest()

	_, err := t.ZeusTestClient.DeployCluster(ctx, cdep)
	t.Require().Nil(err)
}
func (t *ZeusP2PCrawlerTestSuite) TestChartUpload() {
	cd := hardhatClusterDefinition
	_, err := cd.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(err)
}

func (t *ZeusP2PCrawlerTestSuite) TestCreateClusterClass() {
	cd := hardhatClusterDefinition
	gcd := cd.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	err := gcd.CreateClusterClassDefinitions(context.Background(), t.ZeusTestClient)
	t.Require().Nil(err)
}

type ZeusP2PCrawlerTestSuite struct {
	test_suites_base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

func (t *ZeusP2PCrawlerTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	olympus_cookbooks.ChangeToCookbookDir()
}

func TestZeusP2PCrawlerTestSuite(t *testing.T) {
	suite.Run(t, new(ZeusP2PCrawlerTestSuite))
}
