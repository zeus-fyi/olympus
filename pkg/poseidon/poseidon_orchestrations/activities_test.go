package poseidon_orchestrations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	beacon_actions "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacons/actions"
	client_consts "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacons/constants"
	athena_client "github.com/zeus-fyi/olympus/pkg/athena/client"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_buckets"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	test_base "github.com/zeus-fyi/olympus/test"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

var ctx = context.Background()

type PoseidonActivitiesTestSuite struct {
	test_suites_base.TestSuite
}

func (t *PoseidonActivitiesTestSuite) TestPauseClient() {
	cmName := "cm-lighthouse"
	clientName := "lighthouse"
	resp, err := PoseidonSyncWorker.BeaconActionsClient.PauseClient(ctx, cmName, clientName)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *PoseidonActivitiesTestSuite) TestRsyncConsensus() {
	reqHeader := beacon_actions.BeaconKnsReq
	resp, err := PoseidonSyncWorker.UploadViaPortForward(ctx, reqHeader, poseidon_buckets.LighthouseMainnetBucket)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *PoseidonActivitiesTestSuite) TestResumeClient() {
	cmName := "cm-lighthouse"
	clientName := "lighthouse"
	resp, err := PoseidonSyncWorker.BeaconActionsClient.StartClient(ctx, cmName, clientName)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *PoseidonActivitiesTestSuite) TestConsensusSyncStatus() {
	PoseidonSyncWorker.BeaconActionsClient.ConsensusClient = client_consts.Lighthouse
	resp, err := PoseidonSyncWorker.BeaconActionsClient.GetConsensusClientSyncStatus(ctx)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *PoseidonActivitiesTestSuite) TestExecClientStatus() {
	PoseidonSyncWorker.BeaconActionsClient.ExecClient = client_consts.Geth
	resp, err := PoseidonSyncWorker.BeaconActionsClient.GetExecClientSyncStatus(ctx)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *PoseidonActivitiesTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()
	PoseidonSyncWorker.BeaconActionsClient = beacon_actions.NewLocalBeaconActionsClient(tc.Bearer)
	PoseidonSyncWorker.AthenaClient = athena_client.NewLocalAthenaClient(tc.Bearer)
	// points working dir to inside /test
	test_base.ForceDirToTestDirLocation()
}

func TestPoseidonActivitiesTestSuite(t *testing.T) {
	suite.Run(t, new(PoseidonActivitiesTestSuite))
}
