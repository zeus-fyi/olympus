package hestia_quicknode_models

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

var ctx = context.Background()

type QuickNodeProvisioningTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *QuickNodeProvisioningTestSuite) TestInsertProvisionedService() {
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	psBase := hestia_autogen_bases.ProvisionedQuickNodeServices{
		QuickNodeID: uuid.New().String(),
		EndpointID:  uuid.New().String(),
		HttpURL: sql.NullString{
			String: "https://" + uuid.NewString() + ".quiknode.pro/" + uuid.NewString(),
			Valid:  true,
		},
		Network: sql.NullString{
			String: "mainnet",
			Valid:  true,
		},
		Plan:   "lite",
		Active: true,
		OrgID:  s.Tc.ProductionLocalTemporalOrgID,
		WssURL: sql.NullString{
			String: "ws://" + uuid.NewString() + ".quiknode.pro/" + uuid.NewString(),
			Valid:  true,
		},
		Chain: sql.NullString{
			String: "ethereum",
			Valid:  true,
		},
	}

	ps := QuickNodeService{
		ProvisionedQuickNodeServices:                       psBase,
		ProvisionedQuicknodeServicesContractAddressesSlice: nil,
		ProvisionedQuicknodeServicesReferersSlice:          nil,
	}
	err := InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)
}

func (s *QuickNodeProvisioningTestSuite) TestUpdateProvisionedService() {

}

func (s *QuickNodeProvisioningTestSuite) TestDeprovisionedService() {

}

func (s *QuickNodeProvisioningTestSuite) TestDeactivateService() {

}

func TestQuickNodeProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(QuickNodeProvisioningTestSuite))
}
