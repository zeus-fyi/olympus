package ai_platform_service_orchestrations

import (
	"fmt"
	"net/http"
	"time"

	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
)

func (t *ZeusWorkerTestSuite) TestIrisRoutingGroup() {
	artemis_orchestration_auth.Bearer = t.Tc.ProductionLocalTemporalBearerToken
	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID
	var respData []hera_search.SearchResult

	rgName := "test"
	ogr, rerr := iris_models.SelectOrgRoutesByOrgAndGroupName(ctx, ou.OrgID, rgName)
	t.Require().Nil(rerr)
	orgMapGroupRoutes := ogr.Map[ou.OrgID][rgName]
	for _, r := range orgMapGroupRoutes {
		rw := iris_api_requests.NewIrisApiRequestsActivities()
		req := &iris_api_requests.ApiProxyRequest{
			Url:             r.RoutePath,
			PayloadTypeREST: "GET",
			Timeout:         1 * time.Minute,
			StatusCode:      http.StatusOK,
		}

		rr, rrerr := rw.ExtLoadBalancerRequest(ctx, req)
		t.Require().Nil(rrerr)
		sres := hera_search.SearchResult{
			Source: rr.Url,
			Value:  fmt.Sprintf("%s", rr.Response),
			Group:  rgName,
			WebResponse: hera_search.WebResponse{
				Body:       rr.Response,
				RawMessage: rr.RawResponse,
			},
		}
		respData = append(respData, sres)
	}

	t.Require().NotEmpty(respData)
}

func (t *ZeusWorkerTestSuite) TestTgWorkflow() {
	ou := org_users.NewOrgUserWithID(7138983863666903883, 7138958574876245567)

	msgs, err := GetPandoraMessages(ctx, ou, "Zeus")
	t.Require().Nil(err)
	t.Require().NotNil(msgs)

	za := NewZeusAiPlatformActivities()

	for _, msg := range msgs {
		_, err = za.InsertTelegramMessageIfNew(ctx, ou, msg)
		t.Require().Nil(err)
	}
}

func (t *ZeusWorkerTestSuite) TestTokenize() {
	artemis_orchestration_auth.Bearer = t.Tc.ProductionLocalTemporalBearerToken
	tc, err := GetTokenCountEstimate(ctx, "", "The OpenAI Cookbook is a community-driven resource. Whether you're submitting an idea, fixing a typo, adding a new guide, or improving an existing one, your contributions are greatly appreciated!")
	t.Require().Nil(err)
	t.Require().NotZero(tc)
	fmt.Println(tc)
}
