package create_orgs

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type CreateOrgsTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CreateOrgsTestSuite) TestInsertOrg() {
	var ts chronos.Chronos
	s.InitLocalConfigs()
	apps.Pg.InitPG(context.Background(), s.Tc.ProdLocalDbPgconn)
	o := NewCreateNamedOrg("userdemos")
	o.OrgID = ts.UnixTimeStampNow()
	ctx := context.Background()
	err := o.InsertOrg(ctx)
	s.Require().Nil(err)
	fmt.Println(o.OrgID)
}

func TestCreateOrgsTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrgsTestSuite))
}
