package deploy_topology_activities_create_setup

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type DeployTestSuite struct {
	test_suites_base.TestSuite
}

func (s *DeployTestSuite) SetupTest() {
	s.InitLocalConfigs()
}

func (s *DeployTestSuite) TestMakeNodePoolReq() {

	//act := CreateSetupTopologyActivities{}
	//
	//ctx := context.Background()
	//params := base_deploy_params.ClusterSetupRequest{
	//	FreeTrial:  false,
	//	Ou:         org_users.OrgUser{},
	//	CloudCtxNs: zeus_common_types.CloudCtxNs{},
	//	Nodes: autogen_bases.Nodes{
	//		CloudProvider: "do",
	//		Region:        "nyc1",
	//		Slug:          "so1_5-16vcpu-128gb",
	//	},
	//	NodesQuantity: 1,
	//	Disks:         nil,
	//	Cluster:       zeus_templates.Cluster{},
	//	AppTaint:      false,
	//}
	//request, err := act.MakeNodePoolRequest(ctx, params)
	//s.Require().Nil(err)
	//s.Require().NotNil(request)
}

func TestDeployTestSuite(t *testing.T) {
	suite.Run(t, new(DeployTestSuite))
}
