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
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
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

	// todo: add instructions
	//err = os.InsertOrchestrationsWithInstructions(ctx, []byte("test"))
	//s.Require().Nil(err)
	//s.Assert().NotZero(os.OrchestrationID)
	//fmt.Println(os.OrchestrationID)
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
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	orch := artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs{
		Status: "Pending",
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

func (s *OrchestrationsTestSuite) TestSelectActiveInstructions() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ojs, err := SelectActiveOrchestrationsWithInstructions(ctx, s.Tc.ProductionLocalTemporalOrgID, "zeus")
	s.Require().Nil(err)
	s.Assert().NotEmpty(ojs)
}

// SelectOrchestrationsAtCloudCtxNsWithStatus
func TestOrchestrationsTestSuite(t *testing.T) {
	suite.Run(t, new(OrchestrationsTestSuite))
}
