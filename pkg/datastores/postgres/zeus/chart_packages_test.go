package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type ChartPackagesTestSuite struct {
	base.TestSuite
}

var PgTestDB Db

func (s *ChartPackagesTestSuite) TestChartPackagesQuery() {
	packageID := 0
	ctx := context.Background()
	conn := PgTestDB.InitPG(ctx, s.Tc.LocalDbPgconn)
	s.Assert().NotNil(conn)
	defer conn.Close()

	pkg, err := FetchQueryPackage(ctx, packageID)
	s.Require().Nil(err)
	s.Require().NotEmpty(pkg)
}

func TestChartPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(ChartPackagesTestSuite))
}
