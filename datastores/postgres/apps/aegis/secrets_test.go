package aegis_secrets

import (
	"context"
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

var ctx = context.Background()

func (t *AegisSecretsTestSuite) TestInsertSecret() {
	ref := autogen_bases.OrgSecretKeyValReferences{
		SecretID:        0,
		SecretEnvVarRef: "RPC_URL",
		SecretKeyRef:    "rpc",
		SecretNameRef:   "hardhat",
	}
	secRef := autogen_bases.OrgSecretReferences{
		OrgID:      t.Tc.ProductionLocalTemporalOrgID,
		SecretID:   0,
		SecretName: "artemis.ethereum.mainnet.quiknode.txt",
	}
	err := InsertOrgSecretRef(ctx, secRef, ref)
	t.Require().Nil(err)

	refTop := autogen_bases.TopologySystemComponentsSecrets{
		TopologySystemComponentID: 0,
		SecretID:                  0,
	}
	err = InsertOrgSecretTopologyRef(ctx, refTop)
	t.Require().Nil(err)
}

func (t *AegisSecretsTestSuite) TestSelectOrgTopSecretRefs() {
	topId := 0
	orgSecrets, err := SelectOrgSecretRef(ctx, t.Tc.ProductionLocalTemporalOrgID, topId)
	t.Require().Nil(err)
	t.Require().NotNil(orgSecrets)
}
func TestAegisSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(AegisSecretsTestSuite))
}
