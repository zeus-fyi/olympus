package common

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

var ts = chronos.Chronos{}

type ChartSubcomponentParentClassTypesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ChartSubcomponentParentClassTypesTestSuite) TestSeedChartSubcomponentParentClassTypes() {
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("ChartSubcomponentParentClassTypes", "table", "where", 1000, []string{})

	cr := seedHeadlessPodSpecParentClass()
	err := cr.InsertChartSubcomponentParentClassTypes(ctx, q)
	s.Require().Nil(err)
}

func seedHeadlessPodSpecParentClass() ParentClass {
	cr := ParentClass{autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             3,
		ChartSubcomponentParentClassTypeID:   ts.UnixTimeStampNow(),
		ChartSubcomponentParentClassTypeName: "HeadlessPodSpecParent",
	}}
	return cr
}

func TestChartSubcomponentParentClassTypesTestSuite(t *testing.T) {
	suite.Run(t, new(ChartSubcomponentParentClassTypesTestSuite))
}
