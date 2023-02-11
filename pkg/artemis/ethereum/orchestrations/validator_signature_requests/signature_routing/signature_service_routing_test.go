package artemis_validator_signature_service_routing

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	"testing"
)

type ValidatorServiceAuthRoutesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *ValidatorServiceAuthRoutesTestSuite) TestFetchServiceAuthRouteInfo() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	cctx := zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "ephemeral-staking", // set with your own namespace
		Env:           "production",
	}
	svc, err := GetServiceURLs(ctx, cctx)
	s.Require().Nil(err)
	s.Require().NotEmpty(svc)
}

func TestValidatorServiceAuthRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorServiceAuthRoutesTestSuite))
}
