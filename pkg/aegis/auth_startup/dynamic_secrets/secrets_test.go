package dynamic_secrets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	aegis_secrets "github.com/zeus-fyi/olympus/datastores/postgres/apps/aegis"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type DynamicSecretsTestSuite struct {
	s3secrets.S3SecretsManagerTestSuite
}

var ctx = context.Background()

func (t *DynamicSecretsTestSuite) TestSecretLookupAndCreate() {
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)
	keysCfg := auth_keys_config.AuthKeysCfg{
		AgePrivKey:    t.Tc.LocalAgePkey,
		AgePubKey:     t.Tc.LocalAgePubkey,
		SpacesKey:     t.Tc.LocalS3SpacesKey,
		SpacesPrivKey: t.Tc.LocalS3SpacesSecret,
	}
	authCfg := auth_startup.NewDefaultAuthClient(ctx, keysCfg)
	inMemFs := auth_startup.ReadEncryptedSecretsData(ctx, authCfg)
	t.Require().NotEmpty(inMemFs)
	AegisInMemSecrets = inMemFs

	kns := zeus_common_types.CloudCtxNs{
		CloudProvider: "",
		Region:        "",
		Context:       "",
		Namespace:     "test",
		Alias:         "",
		Env:           "",
	}
	topName := "rqhppnzghs"
	sec, err := LookupAndCreateSecret(ctx, t.Tc.ProductionLocalTemporalOrgID, topName, kns)
	t.Require().NoError(err)
	t.Require().NotNil(sec)
	t.Require().NotNil(sec.StringData)
	t.Assert().Equal(t.Tc.MainnetNodeUrl, sec.StringData["rpc"])
}

func (t *DynamicSecretsTestSuite) TestPackageSecret() {
	keysCfg := auth_keys_config.AuthKeysCfg{
		AgePrivKey:    t.Tc.LocalAgePkey,
		AgePubKey:     t.Tc.LocalAgePubkey,
		SpacesKey:     t.Tc.LocalS3SpacesKey,
		SpacesPrivKey: t.Tc.LocalS3SpacesSecret,
	}
	authCfg := auth_startup.NewDefaultAuthClient(ctx, keysCfg)
	inMemFs := auth_startup.ReadEncryptedSecretsData(ctx, authCfg)
	t.Require().NotEmpty(inMemFs)
	AegisInMemSecrets = inMemFs

	sec, err := CreateSecret(ctx, zeus_common_types.CloudCtxNs{
		CloudProvider: "",
		Region:        "",
		Context:       "",
		Namespace:     "test",
		Alias:         "",
		Env:           "",
	}, aegis_secrets.OrgSecretRef{
		OrgSecretReferences: autogen_bases.OrgSecretReferences{
			SecretName: "hardhat",
		},
		OrgSecretKeyValReferencesSlice: []autogen_bases.OrgSecretKeyValReferences{
			{
				SecretEnvVarRef: "RPC_URL",
				SecretNameRef:   "rpc",
				SecretKeyRef:    auth_startup.QuikNodeSecret,
			},
		},
	})
	t.Require().NoError(err)
	t.Require().NotNil(sec)
	t.Require().NotNil(sec.StringData)
	t.Assert().Equal(t.Tc.MainnetNodeUrl, sec.StringData["rpc"])
}
func TestDynamicSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(DynamicSecretsTestSuite))
}
