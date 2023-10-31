package iris_serverless

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
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
	// 7138983863666903883
	//err := iris_redis.IrisRedisClient.ReleaseServerlessRoute(context.Background(), 7138983863666903883, "sessionID", AnvilServerlessRoutingTable)
	//t.Require().NoError(err)

	//route, err := iris_redis.IrisRedisClient.GetNextServerlessRoute(context.Background(), t.Tc.ProductionLocalTemporalOrgID, "sessionID", AnvilServerlessRoutingTable)
	//t.Require().NoError(err)
	//t.Require().NotNil(route)

	//irisSvc := "https://iris.zeus.fyi/v1/router"
	//irisSvc = "http://localhost:8080/v1/router"
	irisSvc := "http://localhost:8888"

	wa := web3_actions.NewWeb3ActionsClient(irisSvc)
	wa.AddDefaultEthereumMainnetTableHeader()
	wa.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	sessionHeader := "X-Anvil-Session-Lock-ID"
	testID := "sessionID"
	wa.Headers[sessionHeader] = testID
	wa.IsAnvilNode = true
	wa.Dial()
	defer wa.Close()

	nodeInfo, err := wa.GetNodeInfo(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(nodeInfo)

	fmt.Println(nodeInfo.ForkConfig.ForkUrl)

}

const (
	NodeURL            = "http://localhost:8545"
	HypnosURL          = "http://localhost:8888"
	LocalProxiedRouter = "http://localhost:8888/node"
)

func (t *IrisOrchestrationsTestSuite) TestAnvilRpcReset() {

	token, err := iris_redis.IrisRedisClient.SetInternalAuthCache(ctx, org_users.OrgUser{
		OrgUsers: hestia_autogen_bases.OrgUsers{
			UserID: t.Tc.ProductionLocalTemporalUserID,
			OrgID:  t.Tc.ProductionLocalTemporalOrgID,
		},
	}, "sessionID", "enterprise", "sessionID")
	t.Require().NoError(err)

	ous, plan, err := iris_redis.IrisRedisClient.GetInternalAuthCacheIfExists(ctx, token)
	t.Require().NoError(err)
	t.Require().NotNil(ous)
	t.Require().NotNil(plan)
	t.Assert().Equal(t.Tc.ProductionLocalTemporalOrgID, ous.OrgID)
	t.Assert().Equal(t.Tc.ProductionLocalTemporalUserID, ous.UserID)
	t.Assert().Equal("enterprise", plan)
	wa := web3_actions.NewWeb3ActionsClient(NodeURL)
	wa.IsAnvilNode = true

	wa.Dial()
	defer wa.Close()

	//wa.AddDefaultEthereumMainnetTableHeader()
	wa.C.Client().SetHeader("X-Anvil-Session-Lock-ID", token)
	wa.C.Client().SetHeader("X-Route-Group", "ethereum-mainnet")

	res, err := wa.SetRpcUrl(ctx, LocalProxiedRouter)
	//err = wa.ResetNetwork(context.Background(), LocalProxiedRouter, 0)
	t.Require().NoError(err)
	fmt.Println(res)
}

func (t *IrisOrchestrationsTestSuite) TestAnvilRpcInfo() {
	na := "http://localhost:8545"
	wa := web3_actions.NewWeb3ActionsClient(na)
	wa.IsAnvilNode = true
	wa.Dial()
	defer wa.Close()

	nodeInfo, err := wa.GetNodeInfo(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(nodeInfo)

	fmt.Println(nodeInfo.ForkConfig.ForkUrl)
}

/*
	wa := web3_actions.NewWeb3ActionsClient(irisSvc)
	wa.AddDefaultEthereumMainnetTableHeader()
	wa.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	sessionHeader := "X-Anvil-Session-Lock-ID"
	testID := "sessionID"
	wa.Headers[sessionHeader] = testID
	wa.IsAnvilNode = true

	//nodeInfo, err := wa.GetNodeInfo(ctx)
	//t.Require().Nil(err)
	//t.Require().NotNil(nodeInfo)
	//
	//fmt.Println(nodeInfo.ForkConfig.ForkUrl)
	na := "http://localhost:8545"

	wa := web3_actions.NewWeb3ActionsClient(na)
	wa.Dial()
	defer wa.Close()

	err := wa.ResetNetwork(context.Background(), na, 0)
	t.Require().NoError(err)

	nodeInfo, err := wa.GetNodeInfo(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(nodeInfo)

	fmt.Println(nodeInfo.ForkConfig.ForkUrl)
*/
// curl --location 'http://anvil-0.anvil.anvil-serverless-4d383226.svc.cluster.local:8545' --header 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}'

func TestIrisOrchestrationsTestSuite(t *testing.T) {
	suite.Run(t, new(IrisOrchestrationsTestSuite))
}
