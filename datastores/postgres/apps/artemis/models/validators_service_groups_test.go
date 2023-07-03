package artemis_validator_service_groups_models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var ctx = context.Background()

type ValidatorServicesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *ValidatorServicesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
}

func (s *ValidatorServicesTestSuite) TestFilterKeysThatExistAlready() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	s.Pg.InitPG(context.Background(), s.Tc.ProdLocalDbPgconn)

	keyOne := hestia_req_types.ValidatorServiceOrgGroup{
		Pubkey:       "0x83b08ac888767416614040a1abc0c0043d2e713454c809813081b31d3e3d59678a67375ca3deebce33bff46c20b98d29",
		FeeRecipient: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
	}
	keyTwo := hestia_req_types.ValidatorServiceOrgGroup{
		Pubkey:       "0x9796d83c0b5d321d2c368bc5583da02b305edc89078c0abe877da76438aa1a224ae71c815ac322b425f247bb582971c2",
		FeeRecipient: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
	}
	pubkeys := hestia_req_types.ValidatorServiceOrgGroupSlice{keyOne, keyTwo}
	reKeys, err := FilterKeysThatExistAlready(ctx, pubkeys)
	s.Require().Nil(err)
	s.Assert().NotEmpty(reKeys)
}
func (s *ValidatorServicesTestSuite) TestSelectInsertUnplacedValidatorsIntoCloudCtxNs() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	vsi := ValidatorServiceCloudCtxNsProtocol{
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		OrgID:             ou.OrgID,
	}

	cctx := zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "ephemeral-staking",
		Env:           "production",
	}
	err := SelectInsertUnplacedValidatorsIntoCloudCtxNs(ctx, vsi, cctx)
	s.Require().Nil(err)
}

func (s *ValidatorServicesTestSuite) TestSelectValidatorsAssignedToCloudCtxNs() {
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

	vals, err := SelectValidatorsAssignedToCloudCtxNs(ctx, vsctx, cctx)
	s.Require().Nil(err)
	s.Assert().NotEmpty(vals)
}

func (s *ValidatorServicesTestSuite) TestInsertValidatorServiceGroup() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	s.Pg.InitPG(context.Background(), s.Tc.ProdLocalDbPgconn)

	vsg := OrgValidatorService{
		GroupName:         "testGroup",
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ServiceURL:        "https://web3signer.zeus.fyi",
		OrgID:             ou.OrgID,
		Enabled:           true,
	}
	keyOne := hestia_req_types.ValidatorServiceOrgGroup{
		Pubkey:       "0x83b08ac888767416614040a1abc0c0043d2e713454c809813081b31d3e3d59678a67375ca3deebce33bff46c20b98d29",
		FeeRecipient: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
	}
	keyTwo := hestia_req_types.ValidatorServiceOrgGroup{
		Pubkey:       "0x83b08ac888767416614040a1abc0c0043d2e713454c809813081b31d3e3d5967da67375ca3deebce33bff46c20b98ddd",
		FeeRecipient: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
	}
	pubkeys := hestia_req_types.ValidatorServiceOrgGroupSlice{keyOne, keyTwo}
	err := InsertVerifiedValidatorsToService(ctx, vsg, pubkeys)
	s.Require().Nil(err)
}

func (s *ValidatorServicesTestSuite) TestSelectValidatorsServiceRoutesAssignedToCloudCtxNs() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	cctx := zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "ephemeral-staking",
		Env:           "production",
	}
	vsi := ValidatorServiceCloudCtxNsProtocol{
		ProtocolNetworkID:     hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ValidatorClientNumber: 0,
	}

	vsGroup, err := SelectValidatorsServiceRoutesAssignedToCloudCtxNs(ctx, vsi, cctx)
	s.Require().Nil(err)
	s.Assert().NotEmpty(vsGroup)
}

func TestValidatorServicesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorServicesTestSuite))
}
