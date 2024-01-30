package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (t *ZeusWorkerTestSuite) TestApiCallRequestTask() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	act := NewZeusAiPlatformActivities()

	rets, err := act.SelectRetrievalTask(ctx, t.Ou, 1706487709357339000)
	t.Require().Nil(err)
	t.Require().NotEmpty(rets)
	ret := rets[0]
	t.Require().Equal(webPlatform, ret.RetrievalPlatform)
	t.Require().NotNil(ret.WebFilters)
	t.Require().NotNil(ret.WebFilters.RoutingGroup)
	r := RouteTask{
		Ou:        t.Ou,
		Retrieval: ret,
		RouteInfo: iris_models.RouteInfo{
			RoutePath: "",
		},
		Payload: nil,
	}
	td, err := act.ApiCallRequestTask(ctx, r)
	t.Require().Nil(err)
	t.Require().NotEmpty(td)
}
