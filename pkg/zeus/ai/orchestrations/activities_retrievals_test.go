package ai_platform_service_orchestrations

import (
	"github.com/labstack/echo/v4"
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
			RoutePath: "https://load-simulator-fe4852b9.zeus.fyi",
		},
		Headers: map[string][]string{
			"X-Sim-Response-Size":   []string{"1"},
			"X-Sim-Response-Format": []string{"json"},
		},
		Payload: echo.Map{
			"foo": "bar",
		},
	}
	td, err := act.ApiCallRequestTask(ctx, r)
	t.Require().Nil(err)
	t.Require().NotEmpty(td)
}
