package artemis_validator_service_groups_models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

var ctx = context.Background()

type ValidatorServicesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
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

func (s *ValidatorServicesTestSuite) TestInsertValidatorServiceGroup() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	vsg := OrgValidatorService{
		GroupName:         "testGroup",
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ServiceURL:        "https://web3signer.zeus.fyi",
		OrgID:             ou.OrgID,
		Enabled:           true,
	}
	keyOne := hestia_req_types.ValidatorServiceOrgGroup{
		Pubkey:       "0x8a7addbf2857a72736205d861169c643545283a74a1ccb71c95dd2c9652acb89de226ca26d60248c4ef9591d7e010288",
		FeeRecipient: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
	}
	keyTwo := hestia_req_types.ValidatorServiceOrgGroup{
		Pubkey:       "0x9a7addbf2857a72736205d861169c643545283a74a1ccb71c95dd2c9652acb89de226ca26d60248c4ef9591d7e010288",
		FeeRecipient: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
	}
	pubkeys := hestia_req_types.ValidatorServiceOrgGroupSlice{keyOne, keyTwo}
	err := InsertVerifiedValidatorsToService(ctx, vsg, pubkeys)
	s.Require().Nil(err)

}

func TestValidatorServicesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorServicesTestSuite))
}
