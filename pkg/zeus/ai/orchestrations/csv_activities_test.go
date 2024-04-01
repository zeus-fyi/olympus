package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
)

func (t *ZeusWorkerTestSuite) TestGetGlobalEntitiesFromRef() {

	ueh := "48123fb7c2365017fec6634c0650dd0ed07c7986796114f9d6d1154f4aaf9acb48123fb7c2365017fec6634c0650dd0ed07c7986796114f9d6d1154f4aaf9acb"
	refs := []artemis_entities.EntitiesFilter{
		{
			Nickname: ueh,
			Platform: "flows",
		},
	}
	ue, err := GetGlobalEntitiesFromRef(ctx, t.Ou, refs)
	t.Require().Nil(err)
	t.Require().NotEmpty(ue)

	for _, v := range ue {
		t.Assert().NotZero(len(v.MdSlice))
	}
}
