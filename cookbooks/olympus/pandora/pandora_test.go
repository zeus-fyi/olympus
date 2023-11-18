package pandora

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	olympus_cookbooks "github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
	zeus_client_ext "github.com/zeus-fyi/zeus/zeus/z_client"
)

type PandoraCookbookTestSuite struct {
	test_suites_base.TestSuite
	ZeusExtTestClient zeus_client_ext.ZeusClient
}

var ctx = context.Background()

func (t *PandoraCookbookTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusExtTestClient = zeus_client_ext.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	olympus_cookbooks.ChangeToCookbookDir()
}

func (t *PandoraCookbookTestSuite) TestCreatePandora() {
	//err := CreatePandora(ctx, t.ZeusExtTestClient, true)
	//t.Require().Nil(err)
	//gcd := pandoraClusterDefinition.BuildClusterDefinitions()
	//t.Assert().NotEmpty(gcd)
	//fmt.Println(gcd)
	//
	//err := gcd.CreateClusterClassDefinitions(context.Background(), t.ZeusExtTestClient)
	//t.Require().Nil(err)
	_, rerr := pandoraClusterDefinition.UploadChartsFromClusterDefinition(ctx, t.ZeusExtTestClient, true)
	t.Require().Nil(rerr)
}

func TestPandoraCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(PandoraCookbookTestSuite))
}
