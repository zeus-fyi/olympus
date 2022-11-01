package read_topology

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type TopologyTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *TopologyTestSuite) TestSelectTopology() {
	tr := NewInfraTopologyReader()

	tr.TopologyID = 6951056435719556916
	tr.OrgID = 1667266332674446258
	tr.UserID = 1667266332670878528
	ctx := context.Background()
	err := tr.SelectTopology(ctx)
	s.Require().Nil(err)

	chart := tr.Chart
	s.Require().Nil(err)
	s.Require().NotEmpty(chart.K8sDeployment)
	s.Require().NotNil(chart.K8sDeployment.Spec.Replicas)
	s.Require().NotEmpty(chart.K8sDeployment.Spec.Template.GetObjectMeta())

	s.Require().NotEmpty(chart.K8sService)
	s.Require().NotEmpty(chart.K8sConfigMap)
	s.Require().NotEmpty(chart.K8sIngress)
	s.Require().NotEmpty(chart.K8sConfigMap.Name)

}

func TestTopologyTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyTestSuite))
}
