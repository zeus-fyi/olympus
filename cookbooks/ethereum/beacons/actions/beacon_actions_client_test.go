package beacon_actions

import (
	"context"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

var ctx = context.Background()

type BeaconActionsTestSuite struct {
	test_suites_base.TestSuite
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
