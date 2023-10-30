package iris_serverless

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
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

	// iris_redis.InitProductionRedisIrisCache(ctx)
	iris_redis.InitLocalTestRedisIrisCache(ctx)
}

func (t *IrisOrchestrationsTestSuite) TestFetchLatestServerlessRoutes() []iris_models.RouteInfo {
	a := NewIrisPlatformActivities()
	routes := a.FetchLatestServerlessRoutes(ctx)
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

func (t *IrisOrchestrationsTestSuite) TestResyncServerlessRoutes() {
	routes := t.TestFetchLatestServerlessRoutes()
	t.Require().NotNil(routes)

	a := NewIrisPlatformActivities()
	err := a.ResyncServerlessRoutes(ctx, routes)
	t.Require().NoError(err)
}

func TestIrisOrchestrationsTestSuite(t *testing.T) {
	suite.Run(t, new(IrisOrchestrationsTestSuite))
}
