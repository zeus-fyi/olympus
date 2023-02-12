package artemis_validator_signature_service_routing

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

func (s *ValidatorServiceAuthRoutesTestSuite) TestFetchServiceAuthRouteGrouping() {
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
	err := GetServiceAuthAndURLs(ctx, cctx)
	s.Require().Nil(err)
}
