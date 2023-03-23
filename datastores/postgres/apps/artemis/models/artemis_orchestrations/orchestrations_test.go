package artemis_orchestrations

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var ctx = context.Background()

type OrchestrationsTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *OrchestrationsTestSuite) TestInsertOrchestrationDefinition() {
	//apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	orch := artemis_autogen_bases.Orchestrations{
		OrgID:             ou.OrgID,
		OrchestrationName: "prysmDataDirDiskWipe",
	}
	os := OrchestrationJob{
		Orchestrations: orch,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: "do",
			Region:        "sfo3",
			Context:       "",
			Namespace:     "",
			Env:           "",
		},
	}

	err := os.InsertOrchestrations(ctx)
	s.Require().Nil(err)
	s.Assert().NotZero(os.OrchestrationID)
	fmt.Println(os.OrchestrationID)
}

func (s *OrchestrationsTestSuite) TestInsertOrchestrationsScheduledToCloudCtxNsUsingName() {
	//apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	orch := artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs{}
	os := OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrchestrationName: "gethDataDirDiskWipe",
		},
		Scheduled: orch,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: "do",
			Region:        "sfo3",
			Context:       "do-sfo3-dev-do-sfo3-zeus",
			Namespace:     "athena-beacon-goerli",
			Env:           "",
		},
	}
	err := os.InsertOrchestrationsScheduledToCloudCtxNsUsingName(ctx)
	s.Require().Nil(err)
	s.Assert().NotZero(os.Scheduled.OrchestrationScheduleID)
}

func (s *OrchestrationsTestSuite) TestUpdateOrchestrationsScheduledToCloudCtxNs() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	orch := artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs{
		OrchestrationID: 1679548290001220864,
		CloudCtxNsID:    1674866203750351872,
		Status:          "Pending",
	}
	os := OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{},
		Scheduled:      orch,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: "do",
			Region:        "sfo3",
			Context:       "do-sfo3-dev-do-sfo3-zeus",
			Namespace:     "athena-beacon-goerli",
			Env:           "",
		},
	}

	err := os.UpdateOrchestrationsScheduledToCloudCtxNs(ctx)
	s.Require().Nil(err)
	s.Assert().NotZero(orch.OrchestrationID)
}

func (s *OrchestrationsTestSuite) TestSelectOrchestrationsAtCloudCtxNsWithStatus() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	orch := artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs{
		Status: "Pending",
	}
	os := OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrchestrationName: "aaa",
		},
		Scheduled: orch,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: "do",
			Region:        "sfo3",
			Context:       "do-sfo3-dev-do-sfo3-zeus",
			Namespace:     "athena-beacon-goerli",
			Env:           "",
		},
	}
	orchTodo, err := os.SelectOrchestrationsAtCloudCtxNsWithStatus(ctx)
	s.Require().Nil(err)
	s.Assert().True(orchTodo)
}

// SelectOrchestrationsAtCloudCtxNsWithStatus
func TestOrchestrationsTestSuite(t *testing.T) {
	suite.Run(t, new(OrchestrationsTestSuite))
}
