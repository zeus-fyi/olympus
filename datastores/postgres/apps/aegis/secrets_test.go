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
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
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

	// test 1671408416567169792
	prod := 1685553743752077000
	refTop := autogen_bases.TopologySystemComponentsSecrets{
		TopologySystemComponentID: prod,
		SecretID:                  secRef.SecretID,
	}
	err = InsertOrgSecretTopologyRef(ctx, refTop)
	t.Require().Nil(err)
}

func (t *AegisSecretsTestSuite) TestInsertSecretTxFetcher() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	ref := autogen_bases.OrgSecretKeyValReferences{
		SecretEnvVarRef: "PG_CONN_STR",
		SecretKeyRef:    auth_startup.PgSecret,
		SecretNameRef:   "postgres-conn-str",
	}
	secRef := autogen_bases.OrgSecretReferences{
		OrgID:      t.Tc.ProductionLocalTemporalOrgID,
		SecretID:   ts.UnixTimeStampNow(),
		SecretName: "postgres-auth",
	}

	err := InsertOrgSecretRef(ctx, secRef, ref)
	t.Require().Nil(err)

	prod := 1684692146539966000
	refTop := autogen_bases.TopologySystemComponentsSecrets{
		TopologySystemComponentID: prod,
		SecretID:                  secRef.SecretID,
	}
	err = InsertOrgSecretTopologyRef(ctx, refTop)
	t.Require().Nil(err)
}

func (t *AegisSecretsTestSuite) TestInsertSecretTxFetcher2() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	//refKey := autogen_bases.OrgSecretKeyValReferences{
	//	SecretEnvVarRef: "DYNAMODB_ACCESS_KEY",
	//	SecretKeyRef:    auth_startup.HydraAccessKeyDynamoDB,
	//	SecretNameRef:   "dynamodb-access-key",
	//}
	refSec := autogen_bases.OrgSecretKeyValReferences{
		SecretEnvVarRef: "DYNAMODB_SECRET_KEY",
		SecretKeyRef:    auth_startup.HydraSecretKeyDynamoDB,
		SecretNameRef:   "dynamodb-secret-key",
	}
	secRef := autogen_bases.OrgSecretReferences{
		OrgID:      t.Tc.ProductionLocalTemporalOrgID,
		SecretID:   ts.UnixTimeStampNow(),
		SecretName: "dynamodb-auth",
	}

	//err := InsertOrgSecretRef(ctx, secRef, refKey)
	//t.Require().Nil(err)
	err := InsertOrgSecretRef(ctx, secRef, refSec)
	t.Require().Nil(err)

	prod := 1684692146539966000
	refTop := autogen_bases.TopologySystemComponentsSecrets{
		TopologySystemComponentID: prod,
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
