package read_topology

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type TopologyTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *TopologyTestSuite) TestSelectTopology() {
	s.InitLocalConfigs()
	//apps.Pg.InitPG(context.Background(), s.Tc.ProdLocalDbPgconn)
	tr := NewInfraTopologyReader()

	tr.TopologyID = 1671004476048440064
	tr.OrgID = s.Tc.ProductionLocalTemporalOrgID
	tr.UserID = s.Tc.ProductionLocalTemporalUserID
	ctx := context.Background()
	err := tr.SelectTopology(ctx)
	s.Require().Nil(err)

	chart := tr.Chart

	// currently dumps to config dir
	p := filepaths.Path{
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
