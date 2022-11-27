package beacon_actions

import (
	"context"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

var ctx = context.Background()

var basePar = zeus_pods_reqs.PodActionRequest{
	TopologyDeployRequest: beaconKnsReq,
	PodName:               "",
	FilterOpts:            nil,
	ClientReq:             nil,
	DeleteOpts:            nil,
}

// set your own topologyID here after uploading a chart workload
var beaconKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1669159384971627008,
	CloudCtxNs: beaconCloudCtxNs,
}

var beaconCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "beacon", // set with your own namespace
	Env:           "dev",
}

type BeaconActionsTestSuite struct {
	base.TestSuite
	BeaconActionsClient
}

func (t *BeaconActionsTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.BeaconActionsClient = NewDefaultBeaconActionsClient(tc.Bearer)
	dir := cookbooks.ChangeToCookbookDir()

	t.BeaconActionsClient.PrintPath.DirIn = path.Join(dir, "/ethereum/beacon/logs")
	t.BeaconActionsClient.PrintPath.DirOut = path.Join(dir, "/ethereum/outputs")
	t.BeaconActionsClient.ConfigPaths.DirIn = "./ethereum/beacon/infra"
	t.BeaconActionsClient.ConfigPaths.DirOut = "./ethereum/outputs"
}

func TestBeaconActionsTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconActionsTestSuite))
}
