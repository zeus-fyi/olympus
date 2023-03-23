package artemis_validator_service_groups_models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"k8s.io/apimachinery/pkg/util/rand"
)

type OrchestrationsTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *OrchestrationsTestSuite) TestInsertOrchestrationDefinition() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	orch := artemis_autogen_bases.Orchestrations{
		OrgID:             ou.OrgID,
		OrchestrationName: "Test" + rand.String(10),
	}
	err := InsertOrchestrations(ctx, &orch)
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

	err := InsertOrchestrationsScheduledToCloudCtxNs(ctx, &orch)
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

	err := UpdateOrchestrationsScheduledToCloudCtxNs(ctx, ou.OrgID, "Test", &orch)
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
	}

	orchTodo, err := SelectOrchestrationsAtCloudCtxNsWithStatus(ctx, ou.OrgID, orch.CloudCtxNsID, "Pending", "Test")
	s.Require().Nil(err)
	s.Assert().True(orchTodo)
}

// SelectOrchestrationsAtCloudCtxNsWithStatus
func TestOrchestrationsTestSuite(t *testing.T) {
	suite.Run(t, new(OrchestrationsTestSuite))
}
