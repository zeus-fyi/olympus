package artemis_orchestrations

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
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
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	orch := artemis_autogen_bases.Orchestrations{
		OrgID:             ou.OrgID,
		OrchestrationName: "GethDataDirDiskWipe",
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
	s.Assert().NotZero(orch.OrchestrationID)
	fmt.Println(orch.OrchestrationID)
}

func (s *OrchestrationsTestSuite) TestInsertOrchestrationsScheduledToCloudCtxNs() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	orch := artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs{
		OrchestrationID: 1679548290001220864,
		CloudCtxNsID:    1674866203750351872,
	}
	os := OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{},
		Scheduled:      orch,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: "do",
			Region:        "sfo3",
			Context:       "",
			Namespace:     "",
			Env:           "",
		},
	}

	err := os.InsertOrchestrationsScheduledToCloudCtxNs(ctx)
	s.Require().Nil(err)
	s.Assert().NotZero(orch.OrchestrationID)
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
			Context:       "",
			Namespace:     "",
			Env:           "",
		},
	}

	err := os.UpdateOrchestrationsScheduledToCloudCtxNs(ctx)
	s.Require().Nil(err)
	s.Assert().NotZero(orch.OrchestrationID)
}

func (s *OrchestrationsTestSuite) TestInsertOrchestrationScheduledToCloudCtxNs() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	orch := artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs{
		OrchestrationID: 1679548290001220864,
		CloudCtxNsID:    1674866203750351872,
		Status:          "Pending",
	}
	os := OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrchestrationID:   0,
			OrgID:             0,
			OrchestrationName: "Test",
		},
		Scheduled: orch,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: "do",
			Region:        "sfo3",
			Context:       "",
			Namespace:     "",
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
