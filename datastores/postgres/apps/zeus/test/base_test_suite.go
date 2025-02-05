package conversions_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
)

type ConversionsTestSuite struct {
	Ts chronos.Chronos
	test_suites.PGTestSuite
	Yr            transformations.YamlFileIO
	TestDirectory string
}

func (s *ConversionsTestSuite) ChangeToTestDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}

func (s *ConversionsTestSuite) SetupTest() {
	s.TestDirectory = s.ChangeToTestDirectory()
	s.Yr = transformations.YamlFileIO{}
	s.InitLocalConfigs()
	s.SetupPGConn()
}

func (s *ConversionsTestSuite) SeedTopology() (int, string) {
	top := topology.NewTopology()
	top.TopologyID = s.Ts.UnixTimeStampNow()
	top.Name = fmt.Sprintf("testTopology_%d", top.TopologyID)
	ctx := context.Background()
	temp := autogen_bases.Topologies{}
	q := sql_query_templates.NewQueryParam("InsertTopology", "topologies", "where", 1000, []string{})
	q.TableName = "topologies"
	q.Columns = temp.GetTableColumns()
	q.Values = []apps.RowValues{apps.RowValues{top.TopologyID, top.Name}}
	_, err := apps.Pg.Exec(ctx, q.InsertSingleElementQuery())
	if err != nil {
		panic(err)
	}
	return top.TopologyID, top.Name
}
