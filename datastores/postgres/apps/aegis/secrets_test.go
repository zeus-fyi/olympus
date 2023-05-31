package aegis_secrets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type AegisSecretsTestSuite struct {
	h hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

var (
	ctx = context.Background()
	ts  chronos.Chronos
)

func (t *AegisSecretsTestSuite) TestInsertSecret() {
	ref := autogen_bases.OrgSecretKeyValReferences{
		SecretEnvVarRef: "RPC_URL",
		SecretKeyRef:    "rpc",
		SecretNameRef:   "hardhat",
	}
	secRef := autogen_bases.OrgSecretReferences{
		OrgID:      t.Tc.ProductionLocalTemporalOrgID,
		SecretID:   ts.UnixTimeStampNow(),
		SecretName: "artemis.ethereum.mainnet.quiknode.txt",
	}
	err := InsertOrgSecretRef(ctx, secRef, ref)
	t.Require().Nil(err)

	refTop := autogen_bases.TopologySystemComponentsSecrets{
		TopologySystemComponentID: 1671408416567169792,
		SecretID:                  secRef.SecretID,
	}
	err = InsertOrgSecretTopologyRef(ctx, refTop)
	t.Require().Nil(err)
}

func (t *AegisSecretsTestSuite) TestSelectOrgTopSecretRefs() {
	topId := 1671408416567169792
	orgSecrets, err := SelectOrgSecretRef(ctx, t.Tc.ProductionLocalTemporalOrgID, topId)
	t.Require().Nil(err)
	t.Require().NotNil(orgSecrets)
}
func TestAegisSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(AegisSecretsTestSuite))
}
