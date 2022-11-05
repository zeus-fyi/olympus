package read_charts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/hera/cookbook/code_templates/models/test"
)

type ChartReaderTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ChartReaderTestSuite) TestSelectQueryName() {
	ctx := context.Background()
	qp := test.CreateTestQueryNameParams()
	qp.QueryName = "SelectSingleChartsResources"
	chart := Chart{}
	chart.ChartPackageID = 1667626930599678541
	qp.CTEQuery.Params = append(qp.CTEQuery.Params, chart.ChartPackageID)
	qp.RawQuery = FetchChartQuery(qp)
	err := chart.SelectSingleChartsResources(ctx, qp)
	s.Require().Nil(err)
	s.Assert().NotEmpty(chart.K8sStatefulSet.Spec.Template.Spec.InitContainers)
	//	s.Require().NotEmpty(chart.K8sDeployment)
	//s.Require().NotNil(chart.K8sDeployment.Spec.Replicas)
	//s.Require().NotEmpty(chart.K8sDeployment.Spec.Template.GetObjectMeta())

	s.Require().NotEmpty(chart.K8sService)
	s.Require().NotEmpty(chart.K8sConfigMap)
	s.Require().NotEmpty(chart.K8sStatefulSet)

}

func TestChartReaderTestSuite(t *testing.T) {
	suite.Run(t, new(ChartReaderTestSuite))
}
