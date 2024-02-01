package hestia_login

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hesta_base_test "github.com/zeus-fyi/olympus/hestia/api/test"
)

type LoginTestSuite struct {
	hesta_base_test.HestiaBaseTestSuite
}

var ctx = context.Background()

func (t *LoginTestSuite) TestLogin() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)
	//t.E.POST(DemoUsersCreateRoute, CreateDemoUserHandler)
	//start := make(chan struct{}, 1)
	//go func() {
	//	close(start)
	//	_ = t.E.Start(":9002")
	//}()
	//<-start
	//defer t.E.Shutdown(ctx)
	//userDemoReq := CreateDemoUserRequest{
	//	Keyname:   "userDemoKey",
	//	Metadata:  `{"name": "username"}`,
	//	ServiceID: create_org_users.EthereumEphemeryServiceID,
	//}
	//resp, err := t.PostRequest(ctx, DemoUsersCreateRoute, userDemoReq)
	//t.Require().Nil(err)
	//fmt.Println(resp)
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
