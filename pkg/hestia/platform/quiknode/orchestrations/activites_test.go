package quicknode_orchestrations

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	apollo_pagerduty "github.com/zeus-fyi/olympus/pkg/apollo/pagerduty"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type HestiaQuickNodeOrchestrationsTestSuite struct {
	test_suites_base.TestSuite
}

func (t *HestiaQuickNodeOrchestrationsTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func (t *HestiaQuickNodeOrchestrationsTestSuite) TestSeed() {
	ai := kronos_helix.Instructions{
		Alerts: kronos_helix.AlertInstructions{
			Severity:  apollo_pagerduty.CRITICAL,
			Message:   "Unable to provision QuickNode services",
			Source:    "HestiaQuickNodeWorkflow",
			Component: "ProvisionWorkflow",
		},
	}
	b, err := json.Marshal(ai)
	t.Require().NoError(err)
	artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplateWithInstructions(
		7138983863666903883, "HestiaQuickNodeWorkflow", "alerts", "temporal", string(b))
}

func (t *HestiaQuickNodeOrchestrationsTestSuite) TestDeduplicateNetworkChain() {
	network := "arb-nova"
	chain := "nova-mainnet"

	groupName := DeduplicateNetworkChain(network, chain)
	t.Require().Equal("arb-nova-mainnet", groupName)

	network = "ethereum-mainnet"
	chain = "mainnet"

	groupName = DeduplicateNetworkChain(network, chain)
	t.Require().Equal("ethereum-mainnet", groupName)

	network = "ethereum"
	chain = "ethereum-mainnet"

	groupName = DeduplicateNetworkChain(network, chain)
	t.Require().Equal("ethereum-mainnet", groupName)

	network = "bsc"
	chain = "bnb-smart-chain"
	groupName = DeduplicateNetworkChain(network, chain)
	t.Require().Equal("bsc-bnb-smart-chain", groupName)
}

func TestHestiaQuickNodeOrchestrationsTestSuite(t *testing.T) {
	suite.Run(t, new(HestiaQuickNodeOrchestrationsTestSuite))
}
