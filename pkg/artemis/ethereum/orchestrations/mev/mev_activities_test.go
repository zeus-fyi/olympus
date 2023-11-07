package artemis_mev_transcations

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
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

func (s *MevActivitiesTestSuite) TestMonitor() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	beaconPath := "https://iris.zeus.fyi/v1/router"
	wc := web3_client.NewWeb3ClientFakeSigner(beaconPath)
	wc.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)
	wc.AddDefaultEthereumMainnetTableHeader()
	aa := NewArtemisMevActivities(wc)

	list, err := aa.MonitorTxStatusReceipts(ctx)
	s.NoError(err)
	s.NotEmpty(list)
}

func (s *MevActivitiesTestSuite) TestWaitForTx() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	beaconPath := "https://iris.zeus.fyi/v1/router"
	wc := web3_client.NewWeb3ClientFakeSigner(beaconPath)
	wc.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)
	wc.AddDefaultEthereumMainnetTableHeader()
	aa := NewArtemisMevActivities(wc)

	ha := accounts.HexToHash("0x7f6b9610da77c92960660af86b6e136668d57ba26f525f2a0033b249ad2c6d4b")
	_, err := aa.WaitForTxReceipt(ctx, ha)
	if strings.Contains(err.Error(), "not found") {
		err = nil
	}

	s.NoError(err)
	//	s.NotEmpty(rx)

}

func TestMevActivitiesTestSuite(t *testing.T) {
	suite.Run(t, new(MevActivitiesTestSuite))
}
