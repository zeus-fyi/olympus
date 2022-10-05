package conversions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TemplateSpecTestSuite struct {
	ChartPackagesTestSuite
}

func (s *TemplateSpecTestSuite) TestChartPackagesQuery() {
	packageID := 0
	ctx := context.Background()
	conn := PgTestDB.InitPG(ctx, s.Tc.LocalDbPgconn)
	s.Assert().NotNil(conn)
	defer conn.Close()

	pkg, err := FetchQueryPackage(ctx, packageID)
	s.Require().Nil(err)
	s.Require().NotEmpty(pkg)
}

func TestTemplateSpecTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateSpecTestSuite))
}
