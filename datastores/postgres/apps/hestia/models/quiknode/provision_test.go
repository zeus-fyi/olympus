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
	qid := uuid.New().String()
	ps := createQnService(qid, uuid.New().String(), "lite")
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

	ps2 := createQnService(qid, uuid.New().String(), "lite")
	err = InsertProvisionedQuickNodeService(ctx, ps2)
	s.Require().Nil(err)
	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 2)
	qns = qnsLookup.EndpointMap[ps.EndpointID]
	s.Require().Len(qns.ProvisionedQuicknodeServicesContractAddresses, 0)
	s.Require().Len(qns.ProvisionedQuicknodeServicesReferers, 1)
	s.Require().Equal("https://google2.com", qns.ProvisionedQuicknodeServicesReferers[0].Referer)
	qns = qnsLookup.EndpointMap[ps2.EndpointID]
	s.Require().Len(qns.ProvisionedQuicknodeServicesContractAddresses, 0)
	s.Require().Len(qns.ProvisionedQuicknodeServicesReferers, 0)

	// UPDATES now with contract addresses
	createQnServiceReferer(&ps2, "https://zoogle.com")
	createQnServiceContractAddr(&ps2, "0x8888")
	err = UpdateProvisionedQuickNodeService(ctx, ps2)
	s.Require().Nil(err)
	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 2)
	qns = qnsLookup.EndpointMap[ps.EndpointID]
	s.Require().Len(qns.ProvisionedQuicknodeServicesContractAddresses, 0)
	s.Require().Len(qns.ProvisionedQuicknodeServicesReferers, 1)
	s.Require().Equal("https://google2.com", qns.ProvisionedQuicknodeServicesReferers[0].Referer)
	qns = qnsLookup.EndpointMap[ps2.EndpointID]
	s.Require().Len(qns.ProvisionedQuicknodeServicesContractAddresses, 1)
	s.Require().Len(qns.ProvisionedQuicknodeServicesReferers, 1)
	s.Require().Equal("https://zoogle.com", qns.ProvisionedQuicknodeServicesReferers[0].Referer)
	s.Require().Equal("0x8888", qns.ProvisionedQuicknodeServicesContractAddresses[0].ContractAddress)
}

func (s *QuickNodeProvisioningTestSuite) TestUpdateProvisionedService() {
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	qnID := uuid.New().String()
	eId := uuid.New().String()
	ps := createQnService(qnID, eId, "standard")
	err := InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)
	qnsLookup, err := SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().Equal("standard", qnsLookup.Plan)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 1)
	s.Require().Equal(eId, qnsLookup.EndpointMap[eId].EndpointID)

	ps.Plan = "standard"
	ps.ProvisionedQuicknodeServicesContractAddresses = hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddressesSlice{
		{
			ContractAddress: "0x01",
		},
	}
	err = UpdateProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)
	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().Equal("standard", qnsLookup.Plan)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 1)
	s.Require().Equal(eId, qnsLookup.EndpointMap[eId].EndpointID)
	s.Require().Len(qnsLookup.EndpointMap[eId].ProvisionedQuicknodeServicesContractAddresses, 1)
	s.Require().Equal("0x01", qnsLookup.EndpointMap[eId].ProvisionedQuicknodeServicesContractAddresses[0].ContractAddress)

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
	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().Equal("performance", qnsLookup.Plan)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 1)
	s.Require().Equal(eId, qnsLookup.EndpointMap[eId].EndpointID)
	s.Require().Len(qnsLookup.EndpointMap[eId].ProvisionedQuicknodeServicesContractAddresses, 1)
	s.Require().Equal("0x02", qnsLookup.EndpointMap[eId].ProvisionedQuicknodeServicesContractAddresses[0].ContractAddress)
	s.Require().Len(qnsLookup.EndpointMap[eId].ProvisionedQuicknodeServicesReferers, 1)
	s.Require().Equal("https://zeus.fyi", qnsLookup.EndpointMap[eId].ProvisionedQuicknodeServicesReferers[0].Referer)
}

func (s *QuickNodeProvisioningTestSuite) TestDeactiveEndpointService() {
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	qid := uuid.New().String()
	eid := uuid.New().String()
	ps := createQnService(qid, eid, "lite")
	createQnServiceReferer(&ps, "https://depro.com")
	createQnServiceContractAddr(&ps, "0x0dddd")
	err := InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)

	eid2 := uuid.New().String()
	ps2 := createQnService(qid, eid2, "lite")
	ps2.ProvisionedQuicknodeServicesReferers = hestia_autogen_bases.ProvisionedQuicknodeServicesReferersSlice{
		{
			Referer: "https://zeus2.fyi",
		},
	}

	err = InsertProvisionedQuickNodeService(ctx, ps2)
	s.Require().Nil(err)
	qnsLookup, err := SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 2)

	_, err = DeactivateProvisionedQuickNodeServiceEndpoint(ctx, qid, eid)
	s.Require().Nil(err)
	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 1)
	s.Require().Equal(eid2, qnsLookup.EndpointMap[eid2].EndpointID)
	s.Require().Len(qnsLookup.EndpointMap[eid2].ProvisionedQuicknodeServicesReferers, 1)
	s.Require().Equal("https://zeus2.fyi", qnsLookup.EndpointMap[eid2].ProvisionedQuicknodeServicesReferers[0].Referer)
	s.Require().Len(qnsLookup.EndpointMap[eid2].ProvisionedQuicknodeServicesContractAddresses, 0)
}

func (s *QuickNodeProvisioningTestSuite) TestDeprovisionedService() {
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	qid := uuid.New().String()
	eid1 := uuid.New().String()
	ps := createQnService(qid, eid1, "performance")
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
	ps2 := createQnService(qid, uuid.New().String(), "performance")
	ps2.ProvisionedQuicknodeServicesReferers = hestia_autogen_bases.ProvisionedQuicknodeServicesReferersSlice{
		{
			Referer: "https://zeus2.fyi",
		},
	}
	err = InsertProvisionedQuickNodeService(ctx, ps2)

	qnsLookup, err := SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().Equal("performance", qnsLookup.Plan)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 2)

	err = DeprovisionQuickNodeServices(ctx, qid)
	s.Require().Nil(err)

	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().NotEmpty(qnsLookup)
	s.Require().Empty(qnsLookup.EndpointMap)
}

func TestQuickNodeProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(QuickNodeProvisioningTestSuite))
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
		IsTest:                       true,
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
