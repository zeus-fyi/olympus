package containers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type PodSpecInsertTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (p *PodSpecInsertTestSuite) TestDummyInsertChartSubcomponentSpecPodTemplateContainers() {
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertChartSubcomponentSpecPodTemplateContainers", "table", "where", 1000, []string{})

	ps := seedHeadlessPodSpec()
	err := ps.InsertChartSubcomponentSpecPodTemplateContainers(ctx, q)
	p.Require().Nil(err)
}

func seedHeadlessPodSpec() PodSpecContainerMetadata {
	cr := PodSpecContainerMetadata{autogen_bases.ChartSubcomponentSpecPodTemplateContainers{
		IsInitContainer:                   false,
		ContainerSortOrder:                0,
		ChartSubcomponentChildClassTypeID: 0,
		ContainerID:                       5819980028125987373,
	}}
	return cr
}

func TestPodSpecInsertTestSuite(t *testing.T) {
	suite.Run(t, new(PodSpecInsertTestSuite))
}
