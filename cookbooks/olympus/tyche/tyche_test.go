package tyche_olympus_cookbook

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
	zeus_client_ext "github.com/zeus-fyi/zeus/zeus/z_client"
)

type TycheCookbookTestSuite struct {
	test_suites_base.TestSuite
	ZeusTestClient    zeus_client.ZeusClient
	ZeusExtTestClient zeus_client_ext.ZeusClient
}

var ctx = context.Background()

func (t *TycheCookbookTestSuite) TestDeploy() {
	_, rerr := TycheClusterDefinition.UploadChartsFromClusterDefinition(ctx, t.ZeusExtTestClient, true)
	t.Require().Nil(rerr)
	//cdep := TycheClusterDefinition.GenerateDeploymentRequest()

	//_, err := t.ZeusExtTestClient.DeployCluster(ctx, cdep)
	//t.Require().Nil(err)
}

func (t *TycheCookbookTestSuite) TestCreateClusterClass() {
	gcd := TycheClusterDefinition.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	err := gcd.CreateClusterClassDefinitions(ctx, t.ZeusExtTestClient)
	t.Require().Nil(err)
}

func (t *TycheCookbookTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	t.ZeusExtTestClient = zeus_client_ext.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	olympus_cookbooks.ChangeToCookbookDir()
}

func TestTycheCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(TycheCookbookTestSuite))
}
