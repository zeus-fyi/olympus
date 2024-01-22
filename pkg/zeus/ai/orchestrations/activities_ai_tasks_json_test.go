package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (t *ZeusWorkerTestSuite) TestJsonModelOutputActivity() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID

	act := NewZeusAiPlatformActivities()

	td, err := act.SelectTaskDefinition(ctx, ou, 1705429002403077000)
	t.Require().Nil(err)
	t.Require().NotEmpty(td)

	for _, task := range td {
		t.Require().NotEmpty(task.Schemas)
	}

	//_, err = act.CreateJsonOutputModelResponse(ctx, ou, hera_openai.OpenAIParams{
	//	Prompt: "This is a test",
	//})
	//t.Require().Nil(err)
}
