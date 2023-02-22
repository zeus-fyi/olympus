package v1hestia

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	hesta_base_test "github.com/zeus-fyi/olympus/hestia/api/test"
)

type HestiaTestSuite struct {
	hesta_base_test.HestiaBaseTestSuite
}

var ctx = context.Background()

func (t *HestiaTestSuite) TestCreateDemoUserEphemery() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)
	t.E.POST(DemoUsersCreateRoute, CreateDemoUserHandler)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9002")
	}()

	<-start
	defer t.E.Shutdown(ctx)
	userDemoReq := CreateDemoUserRequest{
		KeyName:   "userDemoKey",
		Metadata:  `{"name": "username"}`,
		ServiceID: create_org_users.EthereumEphemeryServiceID,
	}
	resp, err := t.PostRequest(ctx, DemoUsersCreateRoute, userDemoReq)
	t.Require().Nil(err)
	fmt.Println(resp)
}

func TestHestiaTestSuite(t *testing.T) {
	suite.Run(t, new(HestiaTestSuite))
}
