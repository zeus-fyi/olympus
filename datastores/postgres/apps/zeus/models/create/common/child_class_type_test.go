package common

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type ChartSubcomponentChildClassTypesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ChartSubcomponentChildClassTypesTestSuite) TestSeedChartSubcomponentChildClassTypes() {
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("ChartSubcomponentChildClassTypes", "table", "where", 1000, []string{})

	cc := seedHeadlessPodSpecClassType()
	err := cc.InsertChartSubcomponentChildClassTypes(ctx, q)
	s.Require().Nil(err)
}

func seedHeadlessPodSpecClassType() ChildClass {
	cc := ChildClass{autogen_bases.ChartSubcomponentChildClassTypes{
		ChartSubcomponentParentClassTypeID:  1666564843324726081,
		ChartSubcomponentChildClassTypeID:   0,
		ChartSubcomponentChildClassTypeName: "HeadlessPodSpecChildType",
	}}
	return cc
}

func TestChartSubcomponentChildClassTypesTestSuite(t *testing.T) {
	suite.Run(t, new(ChartSubcomponentChildClassTypesTestSuite))
}
