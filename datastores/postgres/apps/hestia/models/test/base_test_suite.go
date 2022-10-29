package hestia_test

import (
	"context"
	"os"
	"path"
	"runtime"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/orgs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
)

var PgTestDB apps.Db

type BaseHestiaTestSuite struct {
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

func (s *BaseHestiaTestSuite) SetupTest() {
	s.TestDirectory = ForceDirToCallerLocation()
	s.Yr = transformations.YamlReader{}
	s.InitLocalConfigs()
	s.SetupPGConn()
}

func (s *BaseHestiaTestSuite) NewTestUser() int {
	var ts chronos.Chronos
	ctx := context.Background()

	u := users.NewUser()
	u.UserID = ts.UnixTimeStampNow()

	qu := sql_query_templates.NewQueryParam("NewTestUser", "users", "where", 1000, []string{})
	qu.TableName = u.GetTableName()
	qu.Columns = u.GetTableColumns()
	qu.Values = []apps.RowValues{u.GetRowValues("default")}
	_, err := apps.Pg.Exec(ctx, qu.InsertSingleElementQuery())
	s.Require().Nil(err)
	return u.UserID
}

func (s *BaseHestiaTestSuite) NewTestOrg() int {
	var ts chronos.Chronos
	ctx := context.Background()

	o := orgs.NewOrg()
	o.OrgID = ts.UnixTimeStampNow()

	qo := sql_query_templates.NewQueryParam("NewTestOrg", "orgs", "where", 1000, []string{})
	qo.TableName = o.GetTableName()
	qo.Columns = o.GetTableColumns()
	qo.Values = []apps.RowValues{o.GetRowValues("default")}
	_, err := apps.Pg.Exec(ctx, qo.InsertSingleElementQuery())
	s.Require().Nil(err)
	return o.OrgID
}

func (s *BaseHestiaTestSuite) NewTestOrgAndUser() {
	ctx := context.Background()
	u := users.NewUser()
	u.UserID = s.NewTestUser()
	o := orgs.NewOrg()
	o.OrgID = s.NewTestOrg()

	ou := org_users.NewOrgUser()
	ou.UserID = u.UserID
	ou.OrgID = o.OrgID

	quo := sql_query_templates.NewQueryParam("NewTestOrgUser", "org_users", "where", 1000, []string{})
	quo.TableName = ou.GetTableName()
	quo.Columns = ou.GetTableColumns()
	quo.Values = []apps.RowValues{ou.GetRowValues("default")}

	_, err := apps.Pg.Exec(ctx, quo.InsertSingleElementQuery())
	s.Require().Nil(err)
	return
}
