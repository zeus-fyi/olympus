package ares_test_suite

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
)

type AresTestSuite struct {
	test_suites.TemporalTestSuite
	ProdZeusClient  zeus_client.ZeusClient
	LocalZeusClient zeus_client.ZeusClient
}

func (t *AresTestSuite) SetupTest() {
	t.SetupZeusClients()
}

func (t *AresTestSuite) SetupZeusClients() {
	t.InitLocalConfigs()
	t.ProdZeusClient = zeus_client.NewZeusClient(t.Tc.ProdZeusApiURL, t.Tc.ProductionLocalTemporalBearerToken)
	t.LocalZeusClient = zeus_client.NewZeusClient(t.Tc.LocalZeusApiURL, t.Tc.ProductionLocalTemporalBearerToken)
}

func TestAresTestSuite(t *testing.T) {
	suite.Run(t, new(AresTestSuite))
}
