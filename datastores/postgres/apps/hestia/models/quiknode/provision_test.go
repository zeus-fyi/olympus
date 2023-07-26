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
		ProvisionedQuickNodeServices: psBase,
	}
	err := InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)

	// UPDATES now with referrers
	ps.ProvisionedQuicknodeServicesReferers = hestia_autogen_bases.ProvisionedQuicknodeServicesReferersSlice{
		{
			Referer: "https://google.com",
		},
	}
	ps.ProvisionedQuicknodeServicesContractAddresses = hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddressesSlice{
		{
			ContractAddress: "0x01",
		},
	}
	err = InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)

	ps.ProvisionedQuicknodeServicesReferers = hestia_autogen_bases.ProvisionedQuicknodeServicesReferersSlice{
		{
			Referer: "https://google2.com",
		},
	}
	ps.ProvisionedQuicknodeServicesContractAddresses = nil

	err = InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)
}

func (s *QuickNodeProvisioningTestSuite) TestUpdateProvisionedService() {
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
		ProvisionedQuickNodeServices: psBase,
	}
	err := InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)

	ps.Plan = "standard"
	ps.ProvisionedQuicknodeServicesContractAddresses = hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddressesSlice{
		{
			ContractAddress: "0x01",
		},
	}
	err = UpdateProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)

	ps.Plan = "performance"
	ps.ProvisionedQuicknodeServicesContractAddresses = hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddressesSlice{
		{
			ContractAddress: "0x02",
		},
	}
	ps.ProvisionedQuicknodeServicesReferers = hestia_autogen_bases.ProvisionedQuicknodeServicesReferersSlice{
		{
			Referer: "https://zeus.fyi",
		},
	}
	err = UpdateProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)
}

func (s *QuickNodeProvisioningTestSuite) TestDeprovisionedService() {
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
		Plan:   "standard",
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
		ProvisionedQuickNodeServices: psBase,
	}
	err := InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)

	err = DeactivateProvisionedQuickNodeServiceEndpoint(ctx, psBase.OrgID, psBase.QuickNodeID, psBase.EndpointID)
	s.Require().Nil(err)
}

func (s *QuickNodeProvisioningTestSuite) TestDeactivateService() {
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
		Plan:   "standard",
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
		ProvisionedQuickNodeServices: psBase,
	}
	ps.Plan = "performance"
	ps.ProvisionedQuicknodeServicesContractAddresses = hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddressesSlice{
		{
			ContractAddress: "0x000000001",
		},
	}
	ps.ProvisionedQuicknodeServicesReferers = hestia_autogen_bases.ProvisionedQuicknodeServicesReferersSlice{
		{
			Referer: "https://zeus2.fyi",
		},
	}
	err := InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)

	endpoint := uuid.New().String()
	s.Require().NotEqual(psBase.EndpointID, endpoint)
	psBase2 := hestia_autogen_bases.ProvisionedQuickNodeServices{
		QuickNodeID: psBase.QuickNodeID,
		EndpointID:  endpoint,
		HttpURL: sql.NullString{
			String: "https://" + uuid.NewString() + ".quiknode.pro/" + uuid.NewString(),
			Valid:  true,
		},
		Network: sql.NullString{
			String: "mainnet",
			Valid:  true,
		},
		Plan:   "standard",
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

	ps2 := QuickNodeService{
		ProvisionedQuickNodeServices: psBase2,
	}
	ps2.ProvisionedQuicknodeServicesReferers = hestia_autogen_bases.ProvisionedQuicknodeServicesReferersSlice{
		{
			Referer: "https://zeus2.fyi",
		},
	}
	err = InsertProvisionedQuickNodeService(ctx, ps2)
	s.Require().Nil(err)

	err = DeprovisionQuickNodeServices(ctx, psBase.OrgID, psBase.QuickNodeID)
	s.Require().Nil(err)
}

func TestQuickNodeProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(QuickNodeProvisioningTestSuite))
}
