package ai_platform_service_orchestrations

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (t *ZeusWorkerTestSuite) TestApiCallRequestTask() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	act := NewZeusAiPlatformActivities()

	rets, err := act.SelectRetrievalTask(ctx, t.Ou, 1706767039731058000)
	t.Require().Nil(err)
	t.Require().NotEmpty(rets)
	ret := rets[0]
	t.Require().Equal(apiApproval, ret.RetrievalPlatform)
	t.Require().NotNil(ret.WebFilters)
	t.Require().NotNil(ret.WebFilters.RoutingGroup)
	tmp := "https://api.twitter.com/2/users/"
	r := RouteTask{
		Ou:        t.Ou,
		Retrieval: ret,
		RouteInfo: iris_models.RouteInfo{
			RoutePath: aws.ToString(&tmp),
		},
	}
	td, err := act.ApiCallRequestTask(ctx, r, nil)
	t.Require().Nil(err)
	t.Require().NotEmpty(td)
}
