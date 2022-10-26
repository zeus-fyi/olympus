package read

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

	deploymentValues := Chart{}
	deploymentValues.ChartPackageID = 6399573289093644219
	err := deploymentValues.SelectSingleChartsResources(ctx, qp)
	s.Require().Nil(err)
	s.Require().NotEmpty(deploymentValues)
}

func TestChartReaderTestSuite(t *testing.T) {
	suite.Run(t, new(ChartReaderTestSuite))
}
