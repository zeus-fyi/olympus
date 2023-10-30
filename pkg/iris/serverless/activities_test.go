package iris_serverless

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

/*
FetchLatestServerlessRoutes
*/

var ctx = context.Background()

type IrisOrchestrationsTestSuite struct {
	test_suites_base.TestSuite
}

func (t *IrisOrchestrationsTestSuite) SetupTest() {
	t.InitLocalConfigs()

	iris_redis.InitLocalTestProductionRedisIrisCache(ctx)
	//iris_redis.InitLocalTestRedisIrisCache(ctx)
}

func (t *IrisOrchestrationsTestSuite) TestFetchLatestServerlessRoutes() []iris_models.RouteInfo {
	a := NewIrisPlatformActivities()
	routes, err := a.FetchLatestServerlessRoutes(ctx)
	t.Require().NoError(err)
	t.Require().NotNil(routes)
	for _, route := range routes {
		fmt.Println(route.RoutePath)
	}
	return routes
}

func (t *IrisOrchestrationsTestSuite) TestResyncOnly() {
	a := NewIrisPlatformActivities()
	err := a.ResyncServerlessRoutes(ctx, nil)
	t.Require().NoError(err)
}

// TestAddAndResyncServerlessRoutes can be used to seed the redis cache with routes
func (t *IrisOrchestrationsTestSuite) TestAddAndResyncServerlessRoutes() {
	routes := t.TestFetchLatestServerlessRoutes()
	t.Require().NotNil(routes)

	a := NewIrisPlatformActivities()
	err := a.ResyncServerlessRoutes(ctx, routes)
	t.Require().NoError(err)
}

func (t *IrisOrchestrationsTestSuite) TestAnvilRpc() {
	irisSvc := "https://iris.zeus.fyi/v1/router"
	//irisSvc = "http://localhost:8080/v1/router"
	//irisSvc = "http://localhost:8545"

	wa := web3_actions.NewWeb3ActionsClient(irisSvc)
	wa.AddDefaultEthereumMainnetTableHeader()
	wa.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	sessionHeader := "X-Anvil-Session-Lock-ID"
	testID := "test-1"
	wa.Headers[sessionHeader] = testID
	wa.IsAnvilNode = true
	wa.Dial()
	defer wa.Close()

	nodeInfo, err := wa.GetNodeInfo(ctx)
	t.Require().Nil(err)

	t.Require().NotNil(nodeInfo)
}

func TestIrisOrchestrationsTestSuite(t *testing.T) {
	suite.Run(t, new(IrisOrchestrationsTestSuite))
}
