package hera_openai

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *HeraTestSuite) TestAssistants() {
	ctx := context.Background()

	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	InitHeraOpenAI(s.Tc.OpenAIAuth)

	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	al, err := ListAssistants(ctx, HeraOpenAI)
	s.Require().Nil(err)
	s.Require().NotNil(al)
}
