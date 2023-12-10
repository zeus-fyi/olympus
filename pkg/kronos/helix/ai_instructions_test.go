package kronos_helix

import (
	"encoding/json"
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

// You can change any params for this, it is a template of the other test meant for creating alerts
func (t *KronosWorkerTestSuite) TestInsertAiInstructions() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	inst := Instructions{
		GroupName: "ai",
		Type:      "workflows",
	}

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID
	b, err := json.Marshal(inst)
	t.Require().Nil(err)

	oj := artemis_orchestrations.OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrgID:             ou.OrgID,
			Active:            false,
			GroupName:         inst.GroupName,
			Type:              inst.Type,
			OrchestrationName: fmt.Sprintf("%s-%s", inst.GroupName, inst.Type),
		},
	}
	ojID, err := artemis_orchestrations.InsertOrchestration(ctx, oj, b)
	t.Require().Nil(err)
	t.Assert().NotZero(ojID)
}
