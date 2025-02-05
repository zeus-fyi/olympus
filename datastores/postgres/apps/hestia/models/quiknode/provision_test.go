package hestia_quicknode_models

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

var ctx = context.Background()

type QuickNodeProvisioningTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *QuickNodeProvisioningTestSuite) TestEnterUserExternal() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	//email := "smedwards121@gmail.com'"
	plan := "enterprise"
	uid := 1701381301753642000
	err := InsertIrisUserApiKey(ctx, uid, plan)
	s.Require().Nil(err)
}

func (s *QuickNodeProvisioningTestSuite) TestInsertProvisionedService2() {
	// Initialize PG if necessary
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ps := QuickNodeService{
		IsTest: true,
		ProvisionedQuickNodeServices: hestia_autogen_bases.ProvisionedQuickNodeServices{
			CreatedAt:   time.Time{}, // Assuming current time is not relevant for the test
			UpdatedAt:   time.Time{}, // Assuming current time is not relevant for the test
			QuickNodeID: "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			EndpointID:  "dd939a5c-5dc5-4555-8c6f-213086c2e40d",
			HttpURL:     sql.NullString{String: "https://few-damp-season.quiknode.pro/mktp-835f38df70e74c56511cd301fc1c7454847d8e76/", Valid: true},
			Network:     sql.NullString{String: "mainnet", Valid: true},
			Plan:        "discover",
			Active:      true,                // Assuming you want to set this service as active
			OrgID:       1702329068567473000, // Assuming OrgID is set elsewhere or not relevant for the test
			WssURL:      sql.NullString{String: "wss://few-damp-season.quiknode.pro/mktp-835f38df70e74c56511cd301fc1c7454847d8e76/", Valid: true},
			Chain:       sql.NullString{String: "ethereum", Valid: true},
		},
		ProvisionedQuickNodeServicesContractAddresses: []hestia_autogen_bases.ProvisionedQuickNodeServicesContractAddresses{
			// Populate this array if you have contract addresses to test
		},
		ProvisionedQuickNodeServicesReferrers: []hestia_autogen_bases.ProvisionedQuickNodeServicesReferrers{
			{
				QuickNodeID: "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
				Referer:     "quicknode.com",
			},
			{
				QuickNodeID: "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
				Referer:     "quiknode.io",
			},
		},
	}

	err := UpsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)
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
	s.Require().Len(qns.ProvisionedQuickNodeServicesContractAddresses, 0)
	s.Require().Len(qns.ProvisionedQuickNodeServicesReferrers, 0)
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
	s.Require().Len(qns.ProvisionedQuickNodeServicesContractAddresses, 1)
	s.Require().Len(qns.ProvisionedQuickNodeServicesReferrers, 1)
	s.Require().Equal("https://google1.com", qns.ProvisionedQuickNodeServicesReferrers[0].Referer)
	s.Require().Equal("0x01", qns.ProvisionedQuickNodeServicesContractAddresses[0].ContractAddress)

	ps.ProvisionedQuickNodeServicesReferrers = hestia_autogen_bases.ProvisionedQuickNodeServicesReferersSlice{
		{
			Referer: "https://google2.com",
		},
	}
	ps.ProvisionedQuickNodeServicesContractAddresses = nil

	err = InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)

	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 1)
	qns = qnsLookup.EndpointMap[ps.EndpointID]
	s.Require().Len(qns.ProvisionedQuickNodeServicesContractAddresses, 0)
	s.Require().Len(qns.ProvisionedQuickNodeServicesReferrers, 1)
	s.Require().Equal("https://google2.com", qns.ProvisionedQuickNodeServicesReferrers[0].Referer)

	ps2 := createQnService(qid, uuid.New().String(), "lite")
	err = InsertProvisionedQuickNodeService(ctx, ps2)
	s.Require().Nil(err)
	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 2)
	qns = qnsLookup.EndpointMap[ps.EndpointID]
	s.Require().Len(qns.ProvisionedQuickNodeServicesContractAddresses, 0)
	s.Require().Len(qns.ProvisionedQuickNodeServicesReferrers, 1)
	s.Require().Equal("https://google2.com", qns.ProvisionedQuickNodeServicesReferrers[0].Referer)
	qns = qnsLookup.EndpointMap[ps2.EndpointID]
	s.Require().Len(qns.ProvisionedQuickNodeServicesContractAddresses, 0)
	s.Require().Len(qns.ProvisionedQuickNodeServicesReferrers, 0)

	// UPDATES now with contract addresses
	createQnServiceReferer(&ps2, "https://zoogle.com")
	createQnServiceContractAddr(&ps2, "0x8888")
	err = UpsertProvisionedQuickNodeService(ctx, ps2)
	s.Require().Nil(err)
	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 2)
	qns = qnsLookup.EndpointMap[ps.EndpointID]
	s.Require().Len(qns.ProvisionedQuickNodeServicesContractAddresses, 0)
	s.Require().Len(qns.ProvisionedQuickNodeServicesReferrers, 1)
	s.Require().Equal("https://google2.com", qns.ProvisionedQuickNodeServicesReferrers[0].Referer)
	qns = qnsLookup.EndpointMap[ps2.EndpointID]
	s.Require().Len(qns.ProvisionedQuickNodeServicesContractAddresses, 1)
	s.Require().Len(qns.ProvisionedQuickNodeServicesReferrers, 1)
	s.Require().Equal("https://zoogle.com", qns.ProvisionedQuickNodeServicesReferrers[0].Referer)
	s.Require().Equal("0x8888", qns.ProvisionedQuickNodeServicesContractAddresses[0].ContractAddress)
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
	ps.ProvisionedQuickNodeServicesContractAddresses = hestia_autogen_bases.ProvisionedQuickNodeServicesContractAddressesSlice{
		{
			ContractAddress: "0x01",
		},
	}
	err = UpsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)
	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().Equal("standard", qnsLookup.Plan)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 1)
	s.Require().Equal(eId, qnsLookup.EndpointMap[eId].EndpointID)
	s.Require().Len(qnsLookup.EndpointMap[eId].ProvisionedQuickNodeServicesContractAddresses, 1)
	s.Require().Equal("0x01", qnsLookup.EndpointMap[eId].ProvisionedQuickNodeServicesContractAddresses[0].ContractAddress)

	ps.Plan = "performance"
	ps.ProvisionedQuickNodeServicesContractAddresses = hestia_autogen_bases.ProvisionedQuickNodeServicesContractAddressesSlice{
		{
			ContractAddress: "0x02",
		},
	}
	ps.ProvisionedQuickNodeServicesReferrers = hestia_autogen_bases.ProvisionedQuickNodeServicesReferersSlice{
		{
			Referer: "https://zeus.fyi",
		},
	}
	err = UpsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)
	qnsLookup, err = SelectQuickNodeServicesByQid(ctx, ps.QuickNodeID)
	s.Require().Nil(err)
	s.Require().Equal("performance", qnsLookup.Plan)
	s.Require().NotEmpty(qnsLookup)
	s.Require().NotEmpty(qnsLookup.EndpointMap)
	s.Require().Len(qnsLookup.EndpointMap, 1)
	s.Require().Equal(eId, qnsLookup.EndpointMap[eId].EndpointID)
	s.Require().Len(qnsLookup.EndpointMap[eId].ProvisionedQuickNodeServicesContractAddresses, 1)
	s.Require().Equal("0x02", qnsLookup.EndpointMap[eId].ProvisionedQuickNodeServicesContractAddresses[0].ContractAddress)
	s.Require().Len(qnsLookup.EndpointMap[eId].ProvisionedQuickNodeServicesReferrers, 1)
	s.Require().Equal("https://zeus.fyi", qnsLookup.EndpointMap[eId].ProvisionedQuickNodeServicesReferrers[0].Referer)
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
	ps2.ProvisionedQuickNodeServicesReferrers = hestia_autogen_bases.ProvisionedQuickNodeServicesReferersSlice{
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
	s.Require().Len(qnsLookup.EndpointMap[eid2].ProvisionedQuickNodeServicesReferrers, 1)
	s.Require().Equal("https://zeus2.fyi", qnsLookup.EndpointMap[eid2].ProvisionedQuickNodeServicesReferrers[0].Referer)
	s.Require().Len(qnsLookup.EndpointMap[eid2].ProvisionedQuickNodeServicesContractAddresses, 0)
}

func (s *QuickNodeProvisioningTestSuite) TestDeprovisionedService() {
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	qid := uuid.New().String()
	eid1 := uuid.New().String()
	ps := createQnService(qid, eid1, "performance")
	ps.Plan = "performance"
	ps.ProvisionedQuickNodeServicesContractAddresses = hestia_autogen_bases.ProvisionedQuickNodeServicesContractAddressesSlice{
		{
			ContractAddress: "0x000000001",
		},
	}
	ps.ProvisionedQuickNodeServicesReferrers = hestia_autogen_bases.ProvisionedQuickNodeServicesReferersSlice{
		{
			Referer: "https://zeus2.fyi",
		},
	}
	err := InsertProvisionedQuickNodeService(ctx, ps)
	s.Require().Nil(err)
	ps2 := createQnService(qid, uuid.New().String(), "performance")
	ps2.ProvisionedQuickNodeServicesReferrers = hestia_autogen_bases.ProvisionedQuickNodeServicesReferersSlice{
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
	if qs.ProvisionedQuickNodeServicesContractAddresses == nil {
		qs.ProvisionedQuickNodeServicesContractAddresses = hestia_autogen_bases.ProvisionedQuickNodeServicesContractAddressesSlice{}
	}
	qs.ProvisionedQuickNodeServicesContractAddresses = append(qs.ProvisionedQuickNodeServicesContractAddresses,
		hestia_autogen_bases.ProvisionedQuickNodeServicesContractAddresses{
			ContractAddress: ca,
		},
	)
	return qs
}

func createQnServiceReferer(qs *QuickNodeService, re string) *QuickNodeService {
	if qs.ProvisionedQuickNodeServicesReferrers == nil {
		qs.ProvisionedQuickNodeServicesReferrers = hestia_autogen_bases.ProvisionedQuickNodeServicesReferersSlice{}
	}
	qs.ProvisionedQuickNodeServicesReferrers = append(qs.ProvisionedQuickNodeServicesReferrers,
		hestia_autogen_bases.ProvisionedQuickNodeServicesReferrers{
			Referer: re,
		},
	)
	return qs
}
