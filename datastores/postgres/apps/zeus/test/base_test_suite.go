package conversions_test

import (
	"context"
	"database/sql"
	"os"
	"path"
	"runtime"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
	"k8s.io/apimachinery/pkg/util/rand"
)

var PgTestDB apps.Db

type ConversionsTestSuite struct {
	test_suites.PGTestSuite
	Yr            transformations.YamlReader
	TestDirectory string
}

func ForceDirToCallerLocation() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}

func (s *ConversionsTestSuite) SetupTest() {
	s.TestDirectory = ForceDirToCallerLocation()
	s.Yr = transformations.YamlReader{}
	s.InitLocalConfigs()
	s.SetupPGConn()
}

func (s *ConversionsTestSuite) MockChart() create.Chart {
	ns := sql.NullString{}
	c := create.Chart{ChartPackages: autogen_bases.ChartPackages{
		ChartName:        rand.String(10),
		ChartVersion:     rand.String(10),
		ChartDescription: ns,
	}}
	q := sql_query_templates.NewQueryParam("InsertMockChartForTest", "table", "where", 1000, []string{})
	ctx := context.Background()
	err := c.InsertChart(ctx, q)
	s.Require().Nil(err)

	return c
}
