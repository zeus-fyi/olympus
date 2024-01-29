package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (t *ZeusWorkerTestSuite) TestTriggerActions() {
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)

	//ou := t.Ou
	//act := NewZeusAiPlatformActivities()}

	// TODO test trigger API calls

	// TOOD: update
	// Eval Results: 1706135988888304000, needs to use result id for unique key
}
