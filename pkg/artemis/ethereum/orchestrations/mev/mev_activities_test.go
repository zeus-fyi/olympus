package artemis_mev_transcations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type MevActivitiesTestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (s *MevActivitiesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	artemis_orchestration_auth.Bearer = s.Tc.ProductionLocalTemporalBearerToken
}

func (s *MevActivitiesTestSuite) TestEndServerlessSessionActivity() {
	//urlPath := "https://iris.zeus.fyi/v1/serverless"
	beaconPath := "https://iris.zeus.fyi/v1/router"
	wc := web3_client.NewWeb3ClientFakeSigner(beaconPath)
	wc.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)
	wc.AddDefaultEthereumMainnetTableHeader()
	aa := NewArtemisMevActivities(wc)

	sessionID := "67637b04-a305-4169-9339-3903e0fa2a62"
	err := aa.EndServerlessSession(ctx, sessionID)
	s.NoError(err)
}

func TestMevActivitiesTestSuite(t *testing.T) {
	suite.Run(t, new(MevActivitiesTestSuite))
}
