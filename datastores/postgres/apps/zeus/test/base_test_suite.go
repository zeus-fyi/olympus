package conversions_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
)

var PgTestDB apps.Db

type ConversionsTestSuite struct {
	Ts chronos.Chronos
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

func (s *ConversionsTestSuite) SeedTopology() (int, string) {
	top := topologies.NewTopology()
	top.TopologyID = s.Ts.UnixTimeStampNow()
	top.Name = fmt.Sprintf("testTopology_%s", top.TopologyID)
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertTopology", "topologies", "where", 1000, []string{})
	q.TableName = top.GetTableName()
	q.Columns = top.GetTableColumns()
	q.Values = []apps.RowValues{top.GetRowValues("default")}
	_, err := apps.Pg.Exec(ctx, q.InsertSingleElementQuery())
	s.Require().Nil(err)

	return top.TopologyID, top.Topologies.Name
}
