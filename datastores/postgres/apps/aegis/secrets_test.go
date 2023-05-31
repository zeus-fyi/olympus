package aegis_secrets

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type AegisSecretsTestSuite struct {
	h hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (t *AegisSecretsTestSuite) TestInsertSecret() {
	ref := autogen_bases.OrgSecretKeyValReferences{
		SecretID:        0,
		SecretEnvVarRef: "RPC_URL",
		SecretKeyRef:    "rpc",
		SecretNameRef:   "hardhat",
	}
	fmt.Println(ref)
	secRef := autogen_bases.OrgSecretReferences{
		OrgID:      t.Tc.ProductionLocalTemporalOrgID,
		SecretID:   0,
		SecretName: "artemis.ethereum.mainnet.quiknode.txt",
	}
	fmt.Println(secRef)
	refTop := autogen_bases.TopologySystemComponentsSecrets{
		TopologySystemComponentID: 0,
		SecretID:                  0,
	}
	fmt.Println(refTop)
}

func TestAegisSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(AegisSecretsTestSuite))
}
