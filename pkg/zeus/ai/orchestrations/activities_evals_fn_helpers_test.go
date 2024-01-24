package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (t *ZeusWorkerTestSuite) TestTransformJSONToEvalScoredMetrics() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID

	act := NewZeusAiPlatformActivities()
	evalFns, err := act.EvalLookup(ctx, ou, 1704066747085827000)
	t.Require().Nil(err)
	t.Require().NotEmpty(evalFns)

	for _, evalFnWithMetrics := range evalFns {
		fmt.Println(evalFnWithMetrics.EvalType)
		resp, rerr := EvalModelScoredJsonOutput(ctx, &evalFnWithMetrics)
		t.Require().Nil(rerr)
		t.Require().NotNil(resp)
	}
}
