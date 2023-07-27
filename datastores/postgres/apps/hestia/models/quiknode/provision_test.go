package hestia_quicknode_models

import (
	"context"
	"database/sql"
	"fmt"
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

func createQnService(qnId, eId, plan string) QuickNodeService {
	psBase := hestia_autogen_bases.ProvisionedQuickNodeServices{
		QuickNodeID: qnId,
		EndpointID:  eId,
		HttpURL: sql.NullString{
			String: "https://" + uuid.NewString() + ".quiknode.pro/" + uuid.NewString(),
			Valid:  true,
		},
		Network: sql.NullString{
			String: "mainnet",
			Valid:  true,
		},
		Plan:   plan,
		Active: true,
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
	return ps
}

func createQnServiceContractAddr(qs *QuickNodeService, ca string) *QuickNodeService {
	if qs.ProvisionedQuicknodeServicesContractAddresses == nil {
		qs.ProvisionedQuicknodeServicesContractAddresses = hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddressesSlice{}
	}
	qs.ProvisionedQuicknodeServicesContractAddresses = append(qs.ProvisionedQuicknodeServicesContractAddresses,
		hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses{
			ContractAddress: ca,
		},
	)
	return qs
}

func createQnServiceReferer(qs *QuickNodeService, re string) *QuickNodeService {
	if qs.ProvisionedQuicknodeServicesReferers == nil {
		qs.ProvisionedQuicknodeServicesReferers = hestia_autogen_bases.ProvisionedQuicknodeServicesReferersSlice{}
	}
	qs.ProvisionedQuicknodeServicesReferers = append(qs.ProvisionedQuicknodeServicesReferers,
		hestia_autogen_bases.ProvisionedQuicknodeServicesReferers{
			Referer: re,
		},
	)
	return qs
}

func (s *QuickNodeProvisioningTestSuite) TestInsertProvisionedService() {
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	ps := createQnService(uuid.New().String(), uuid.New().String(), "lite")
	err := InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)

	qnsLookup, err := SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().Equal("lite", qnsLookup.Plan)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 1)

	qns := qnsLookup.EndpointMap[ps.EndpointID]
	s.Require().Len(qns.ProvisionedQuicknodeServicesContractAddresses, 0)
	s.Require().Len(qns.ProvisionedQuicknodeServicesReferers, 0)
	// UPDATES now with referrers
	createQnServiceReferer(&ps, "https://google1.com")
	createQnServiceContractAddr(&ps, "0x01")

	err = InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)
	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 1)
	qns = qnsLookup.EndpointMap[ps.EndpointID]
	s.Require().Len(qns.ProvisionedQuicknodeServicesContractAddresses, 1)
	s.Require().Len(qns.ProvisionedQuicknodeServicesReferers, 1)
	s.Require().Equal("https://google1.com", qns.ProvisionedQuicknodeServicesReferers[0].Referer)
	s.Require().Equal("0x01", qns.ProvisionedQuicknodeServicesContractAddresses[0].ContractAddress)

	ps.ProvisionedQuicknodeServicesReferers = hestia_autogen_bases.ProvisionedQuicknodeServicesReferersSlice{
		{
			Referer: "https://google2.com",
		},
	}
	ps.ProvisionedQuicknodeServicesContractAddresses = nil

	err = InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)

	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 1)
	qns = qnsLookup.EndpointMap[ps.EndpointID]
	s.Require().Len(qns.ProvisionedQuicknodeServicesContractAddresses, 0)
	s.Require().Len(qns.ProvisionedQuicknodeServicesReferers, 1)
	s.Require().Equal("https://google2.com", qns.ProvisionedQuicknodeServicesReferers[0].Referer)
}

func (s *QuickNodeProvisioningTestSuite) TestUpdateProvisionedService() {
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	qnID := uuid.New().String()
	psBase := hestia_autogen_bases.ProvisionedQuickNodeServices{
		QuickNodeID: qnID,
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

	fmt.Println("psBase.QuickNodeID", psBase.QuickNodeID)
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

	err = DeactivateProvisionedQuickNodeServiceEndpoint(ctx, psBase.QuickNodeID, psBase.EndpointID)
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

	err = DeprovisionQuickNodeServices(ctx, psBase.QuickNodeID)
	s.Require().Nil(err)
}

func TestQuickNodeProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(QuickNodeProvisioningTestSuite))
}
