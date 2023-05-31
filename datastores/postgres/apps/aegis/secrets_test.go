package aegis_secrets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
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

func (t *AegisSecretsTestSuite) TestDoesExist() {
	topName := "rqhppnzghs"
	exists, err := DoesOrgSecretExistForTopology(ctx, t.Tc.ProductionLocalTemporalOrgID, topName)
	t.Require().NoError(err)
	t.Assert().True(exists)

	topName = "asdf"
	exists, err = DoesOrgSecretExistForTopology(ctx, t.Tc.ProductionLocalTemporalOrgID, topName)
	t.Require().NoError(err)
	t.Assert().False(exists)
}

func (t *AegisSecretsTestSuite) TestInsertSecret() {

	ref := autogen_bases.OrgSecretKeyValReferences{
		SecretEnvVarRef: "RPC_URL",
		SecretKeyRef:    auth_startup.QuikNodeSecret,
		SecretNameRef:   "rpc",
	}
	secRef := autogen_bases.OrgSecretReferences{
		OrgID:      t.Tc.ProductionLocalTemporalOrgID,
		SecretID:   ts.UnixTimeStampNow(),
		SecretName: "hardhat",
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
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)
	topName := "rqhppnzghs"
	orgSecrets, err := SelectOrgSecretRef(ctx, t.Tc.ProductionLocalTemporalOrgID, topName)
	t.Require().Nil(err)
	t.Require().NotNil(orgSecrets)
}
func TestAegisSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(AegisSecretsTestSuite))
}
