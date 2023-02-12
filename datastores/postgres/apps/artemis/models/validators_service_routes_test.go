package artemis_validator_service_groups_models

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type ValidatorServiceRoutesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *ValidatorServiceRoutesTestSuite) TestFetchServiceRouteInfo() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	vsctx := ValidatorServiceCloudCtxNsProtocol{
		ProtocolNetworkID:     hestia_req_types.EthereumEphemeryProtocolNetworkID,
		OrgID:                 ou.OrgID,
		ValidatorClientNumber: 0,
	}
	cctx := zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "ephemeral-staking",
		Env:           "production",
	}
	vsr, err := SelectValidatorsServiceRoutesAssignedToCloudCtxNs(ctx, vsctx, cctx)
	s.Require().Nil(err)
	s.Require().NotNil(vsr.PubkeyToGroupName)
	s.Require().NotEmpty(vsr.PubkeyToGroupName)
}

func TestValidatorServiceRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorServiceRoutesTestSuite))
}
