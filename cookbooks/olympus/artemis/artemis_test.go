package artemis_cookbook

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/internal_reqs"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type ArtemisCookbookTestSuite struct {
	test_suites_base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

func (t *ArtemisCookbookTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	olympus_cookbooks.ChangeToCookbookDir()
}

func (t *ArtemisCookbookTestSuite) TestDeploy() {
	resp, err := t.ZeusTestClient.Deploy(ctx, ArtemisDeployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
	t.TestArtemisSecretsCopy()
}

func (t *ArtemisCookbookTestSuite) TestUploadCharts() {
	resp, err := t.ZeusTestClient.UploadChart(ctx, artemisChartPath, uploadChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)
	ArtemisDeployKnsReq.TopologyID = resp.TopologyID
	tar := zeus_req_types.TopologyRequest{TopologyID: ArtemisDeployKnsReq.TopologyID}
	chartResp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)
	err = chartResp.PrintWorkload(artemisChartPath)
	t.Require().Nil(err)
}

func (t *ArtemisCookbookTestSuite) TestArtemisSecretsCopy() {
	s1 := "spaces-auth"
	s2 := "spaces-key"
	s3 := "age-auth"
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
		ToKns: kns.TopologyKubeCtxNs{
			TopologyID: 0,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "do",
				Region:        "sfo3",
				Context:       "do-sfo3-dev-do-sfo3-zeus",
				Namespace:     "artemis",
				Env:           "dev",
			},
		},
	}
	err := t.ZeusTestClient.CopySecretsFromToNamespace(ctx, req)
	t.Require().Nil(err)
}

func TestArtemisCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisCookbookTestSuite))
}
