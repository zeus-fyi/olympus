package poseidon_orchestrations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	athena_client "github.com/zeus-fyi/olympus/pkg/athena/client"
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_buckets"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	test_base "github.com/zeus-fyi/olympus/test"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
	beacon_actions "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons/actions"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

var ctx = context.Background()

type PoseidonActivitiesTestSuite struct {
	test_suites_base.TestSuite
}

func (t *PoseidonActivitiesTestSuite) TestPauseClient() {
	cmName := "cm-geth"
	clientName := "geth"
	resp, err := PoseidonSyncActivitiesOrchestrator.BeaconActionsClient.PauseClient(ctx, cmName, clientName)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *PoseidonActivitiesTestSuite) TestRsyncConsensus() {
	reqHeader := zeus_req_types.TopologyDeployRequest{
		TopologyID:                      0,
		RequestChoreographySecretDeploy: false,
	}
	resp, err := PoseidonSyncActivitiesOrchestrator.UploadViaPortForward(ctx, reqHeader, poseidon_buckets.LighthouseMainnetBucket)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *PoseidonActivitiesTestSuite) TestRsyncExec() {
	reqHeader := zeus_req_types.TopologyDeployRequest{
		TopologyID:                      0,
		RequestChoreographySecretDeploy: false,
	}
	resp, err := PoseidonSyncActivitiesOrchestrator.UploadViaPortForward(ctx, reqHeader, poseidon_buckets.GethMainnetBucket)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *PoseidonActivitiesTestSuite) TestResumeClient() {
	cmName := "cm-geth"
	clientName := "geth"
	resp, err := PoseidonSyncActivitiesOrchestrator.BeaconActionsClient.StartClient(ctx, cmName, clientName)
	t.Assert().Nil(err)
	t.Assert().NotEmpty(resp)
}

//func (t *PoseidonActivitiesTestSuite) TestConsensusSyncStatus() {
//	PoseidonSyncActivitiesOrchestrator.BeaconActionsClient.ConsensusClient = client_consts.Lighthouse
//	resp, err := PoseidonSyncActivitiesOrchestrator.IsConsensusClientSynced(ctx)
//	t.Assert().Nil(err)
//	t.Assert().NotEmpty(resp)
//}
//
//func (t *PoseidonActivitiesTestSuite) TestExecClientStatus() {
//	PoseidonSyncActivitiesOrchestrator.BeaconActionsClient.ExecClient = client_consts.Geth
//	resp, err := PoseidonSyncActivitiesOrchestrator.IsExecClientSynced(ctx)
//	t.Assert().Nil(err)
//	t.Assert().NotEmpty(resp)
//}

func (t *PoseidonActivitiesTestSuite) SetupTest() {
	tc := api_configs.InitLocalTestConfigs()
	PoseidonSyncActivitiesOrchestrator = NewPoseidonSyncActivity(beacon_actions.NewDefaultBeaconActionsClient(tc.Bearer, kCtxNsHeader), athena_client.NewLocalAthenaClient(tc.Bearer))
	test_base.ForceDirToTestDirLocation()
}

func TestPoseidonActivitiesTestSuite(t *testing.T) {
	suite.Run(t, new(PoseidonActivitiesTestSuite))
}
