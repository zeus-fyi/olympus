package artemis_orchestrations

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestInsertAiOrchestrations() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	res, err := SelectWorkflowTemplates(ctx, ou)
	s.Require().Nil(err)
	s.Require().NotEmpty(res)
	res2, err := GetAiOrchestrationParams(ctx, ou, 0, 0, res)
	s.Require().Nil(err)
	s.Require().NotEmpty(res2)

	for _, wf := range res2 {

		for _, task := range wf.WorkflowTasks {
			fmt.Println(task.AnalysisCycleCount)
		}

	}

}
func (s *OrchestrationsTestSuite) TestCalculateTimeWindow() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	now := int(time.Now().Unix())

	// current cycle == cycleEnd
	// cycleStart = cycleEnd - taskCycleCount
	tw := CalculateTimeWindow(now, 0, 3, time.Minute*5)
	fmt.Println(tw.UnixStartTime, tw.UnixEndTime, tw.UnixEndTime-tw.UnixStartTime)
	fmt.Println(tw.Start, tw.End)
}

func (s *OrchestrationsTestSuite) TestName() {

}
