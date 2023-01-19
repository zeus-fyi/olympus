package artemis_validator_service_groups_models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

var ctx = context.Background()

type ValidatorServicesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *ValidatorServicesTestSuite) TestSelectUnplacedValidators() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
}

func (s *ValidatorServicesTestSuite) TestSelectInsertUnplacedValidatorsIntoCloudCtxNs() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
}

func (s *ValidatorServicesTestSuite) TestSelectValidatorsAssignedToCloudCtxNs() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
}

func TestValidatorServicesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorServicesTestSuite))
}
