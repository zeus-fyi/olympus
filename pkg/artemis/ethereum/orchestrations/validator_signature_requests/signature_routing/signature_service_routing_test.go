package artemis_validator_signature_service_routing

import (
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"testing"
)

type ValidatorServiceAuthRoutesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *ValidatorServiceAuthRoutesTestSuite) TestFetchServiceAuthRouteInfo() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	// TODO
}

func TestValidatorServiceAuthRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorServiceAuthRoutesTestSuite))
}
