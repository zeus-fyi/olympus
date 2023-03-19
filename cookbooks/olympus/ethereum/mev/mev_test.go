package olympus_ethereum_mev_cookbooks

import (
	"testing"

	"github.com/stretchr/testify/suite"
	olympus_cookbooks "github.com/zeus-fyi/olympus/cookbooks"
	olympus_hydra_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/hydra"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/zeus/test/test_suites"
)

type MevCookbookTestSuite struct {
	test_suites.BaseTestSuite
	ZeusTestClient zeus_client.ZeusClient
}

// This assumes a pre-existing cluster class called "hydraGoerli", since
// the cluster already exists, you will run into an error if you try to create
// the cluster class again, so you will use the below functions to add the
// component base and skeleton base to the existing cluster class
func (t *MevCookbookTestSuite) TestCreateClusterBase() {
	basesInsert := []string{"mev"}
	cc := zeus_req_types.TopologyCreateOrAddComponentBasesToClassesRequest{
		ClusterClassName:   olympus_hydra_cookbooks.HydraGoerli,
		ComponentBaseNames: basesInsert,
	}
	_, err := t.ZeusTestClient.AddComponentBasesToClass(ctx, cc)
	t.Require().Nil(err)
}

func (t *MevCookbookTestSuite) TestCreateClusterSkeletonBases() {
	cc := zeus_req_types.TopologyCreateOrAddSkeletonBasesToClassesRequest{
		ClusterClassName:  olympus_hydra_cookbooks.HydraGoerli,
		ComponentBaseName: "mev",
		SkeletonBaseNames: []string{"mevBoost"},
	}
	_, err := t.ZeusTestClient.AddSkeletonBasesToClass(ctx, cc)
	t.Require().Nil(err)
}

func (t *MevCookbookTestSuite) SetupTest() {
	olympus_cookbooks.ChangeToCookbookDir()

	tc := api_configs.InitLocalTestConfigs()
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
}

func TestMevCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(MevCookbookTestSuite))
}
