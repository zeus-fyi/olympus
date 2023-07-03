package olympus_hydra_validators_cookbooks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	olympus_cookbooks "github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	ethereum_web3signer_actions "github.com/zeus-fyi/zeus/cookbooks/ethereum/web3signers/actions"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

const HydraAddress = "http://zeus-hydra:9000"

var ValidatorCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "ephemeral-staking", // set with your own namespace
	Env:           "production",
}

type ValidatorsTestSuite struct {
	test_suites_base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

var ctx = context.Background()

func (t *ValidatorsTestSuite) TestImportValidatorsToSim() {
	req, err := GetValidatorsAndPrepareRemoteSignEmaulation(ctx, t.Tc.EphemeralNodeUrl, HydraAddress)
	t.Require().Nil(err)

	w3 := ethereum_web3signer_actions.Web3SignerActionsClient{ZeusClient: t.ZeusTestClient}
	resp, err := w3.GetLighthouseAuth(ctx, ValidatorCloudCtxNs)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	token := resp.GetAnyValue()
	r, err := w3.EnableWeb3SignerLighthouse(ctx, ValidatorCloudCtxNs, req.Slice, string(token))
	t.Require().Nil(err)
	t.Assert().NotEmpty(r)
}

func (t *ValidatorsTestSuite) SetupTest() {
	olympus_cookbooks.ChangeToCookbookDir()

	t.InitLocalConfigs()

	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(t.Tc.ProductionLocalTemporalBearerToken)
}

func TestValidatorsTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorsTestSuite))
}
