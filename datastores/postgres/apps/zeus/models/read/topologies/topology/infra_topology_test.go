package read_topology

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type TopologyTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *TopologyTestSuite) TestSelectTopology() {
	s.InitLocalConfigs()
	tr := NewInfraTopologyReader()

	tr.TopologyID = 1668065557558818048
	tr.OrgID = 1668065557527643728
	tr.UserID = 1668065557509163089
	ctx := context.Background()
	err := tr.SelectTopology(ctx)
	s.Require().Nil(err)

	chart := tr.Chart

	// currently dumps to config dir
	p := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./",
		FnIn:        "",
		FnOut:       "deployment.yaml",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	b, err := json.Marshal(chart.Deployment.K8sDeployment)
	s.Require().Nil(err)

	err = s.Yr.WriteYamlConfig(p, b)
	s.Require().Nil(err)

	s.Require().NotEmpty(chart.K8sDeployment)
	s.Require().NotNil(chart.K8sDeployment.Spec.Replicas)
	s.Require().NotEmpty(chart.K8sDeployment.Spec.Template.GetObjectMeta())

	//s.Require().NotEmpty(chart.K8sService)
	//s.Require().NotEmpty(chart.K8sConfigMap)
	//s.Require().NotEmpty(chart.K8sConfigMap.Name)

}

func TestTopologyTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyTestSuite))
}
