package hestia_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/orgs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/users"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	"k8s.io/apimachinery/pkg/util/rand"
)

var PgTestDB apps.Db

type BaseHestiaTestSuite struct {
	test_suites.PGTestSuite
	Yr            transformations.YamlFileIO
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

func (b *BaseHestiaTestSuite) SetupTest() {
	b.TestDirectory = ForceDirToCallerLocation()
	b.Yr = transformations.YamlFileIO{}
	b.InitLocalConfigs()
	b.SetupPGConn()
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: b.Tc.AwsAccessKeySecretManager,
		SecretKey: b.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(context.Background(), auth)
}

func (b *BaseHestiaTestSuite) NewTestUser() int {
	var ts chronos.Chronos
	ctx := context.Background()

	u := users.NewUser()
	u.UserID = ts.UnixTimeStampNow()

	qu := sql_query_templates.NewQueryParam("NewTestUser", "users", "where", 1000, []string{})
	qu.TableName = u.GetTableName()
	qu.Columns = u.GetTableColumns()
	qu.Values = []apps.RowValues{u.GetRowValues("default")}
	_, err := apps.Pg.Exec(ctx, qu.InsertSingleElementQuery())
	if err != nil {
		panic(err)
	}
	return u.UserID
}

func (b *BaseHestiaTestSuite) NewTestOrg() int {
	var ts chronos.Chronos
	ctx := context.Background()

	o := orgs.NewOrg()
	o.OrgID = ts.UnixTimeStampNow()
	o.Name = "test_ord" + rand.String(10)
	qo := sql_query_templates.NewQueryParam("NewTestOrg", "orgs", "where", 1000, []string{})
	qo.TableName = o.GetTableName()
	qo.Columns = o.GetTableColumns()
	qo.Values = []apps.RowValues{o.GetRowValues("default")}
	_, err := apps.Pg.Exec(ctx, qo.InsertSingleElementQuery())
	if err != nil {
		panic(err)
	}
	return o.OrgID
}

func (b *BaseHestiaTestSuite) NewTestOrgAndUser() (int, int) {
	ctx := context.Background()
	u := users.NewUser()
	u.UserID = b.NewTestUser()
	o := orgs.NewOrg()
	o.OrgID = b.NewTestOrg()

	ou := org_users.NewOrgUser()
	ou.OrgID = o.OrgID
	ou.UserID = u.UserID

	quo := sql_query_templates.NewQueryParam("NewTestOrgUser", "org_users", "where", 1000, []string{})
	quo.TableName = ou.GetTableName()
	quo.Columns = ou.GetTableColumns()
	quo.Values = []apps.RowValues{ou.GetRowValues("default")}

	_, err := apps.Pg.Exec(ctx, quo.InsertSingleElementQuery())
	if err != nil {
		panic(err)
	}
	return ou.OrgID, ou.UserID
}

func (b *BaseHestiaTestSuite) NewTestOrgAndUserWithBearer() (int, int, string) {
	ctx := context.Background()
	oID, uID := b.NewTestOrgAndUser()
	var ts chronos.Chronos

	nk := autogen_bases.UsersKeys{
		UserID:            uID,
		PublicKeyName:     fmt.Sprintf("bearerkey_%d", ts.UnixTimeStampNow()),
		PublicKeyVerified: true,
		PublicKeyTypeID:   keys.BearerKeyTypeID,
		CreatedAt:         time.Now().UTC(),
		PublicKey:         rand.String(15),
	}
	nk.PublicKeyName = "test_key"
	nk.PublicKeyTypeID = keys.EcdsaKeyTypeID
	nk.CreatedAt = time.Now()
	q := sql_query_templates.NewQueryParam("InsertUserKey", "user_keys", "where", 1000, []string{})
	q.TableName = nk.GetTableName()
	q.Columns = nk.GetTableColumns()
	q.Values = []apps.RowValues{nk.GetRowValues("default")}

	_, err := apps.Pg.Exec(ctx, q.InsertSingleElementQuery())
	if err != nil {
		panic(err)
	}
	return oID, uID, nk.PublicKey
}
