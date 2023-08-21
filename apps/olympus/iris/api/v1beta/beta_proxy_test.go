package v1Beta_iris

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_base_test "github.com/zeus-fyi/olympus/iris/api/test"
)

type ProxyTestSuite struct {
	iris_base_test.IrisBaseTestSuite
}

var ctx = context.Background()

func (t *ProxyTestSuite) TestRoundRobin() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)
	t.E.POST("/v1beta/internal/router/group", InternalRoundRobinRequestHandler)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":8080")
	}()
	<-start
	defer t.E.Shutdown(ctx)

	// TODO: out of date now, not using query params to specify route group
	//iris_round_robin.InitRoutingTables(ctx)
	//m := make(map[string]interface{})
	//m["method"] = "eth_blockNumber"
	//resp, err := t.PostRequest(ctx, "/v1beta/internal/router/group?routeGroup=quiknode-mainnet", m)
	//t.Require().Nil(err)
	//t.Require().NotNil(resp)
}

func TestProxyTestSuite(t *testing.T) {
	suite.Run(t, new(ProxyTestSuite))
}
