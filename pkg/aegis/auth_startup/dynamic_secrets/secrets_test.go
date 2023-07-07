package dynamic_secrets

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	aegis_secrets "github.com/zeus-fyi/olympus/datastores/postgres/apps/aegis"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type DynamicSecretsTestSuite struct {
	s3secrets.S3SecretsManagerTestSuite
}

var ctx = context.Background()

func (t *DynamicSecretsTestSuite) TestGenAndSave() {
	t.InitLocalConfigs()
	pubKey := t.Tc.LocalAgePubkey
	privKey := t.Tc.LocalAgePkey
	age := encryption.NewAge(privKey, pubKey)
	now := time.Now()
	err := SaveAddress(ctx, 100000, t.S3, age)
	t.Require().NoError(err)
	fmt.Println("search time", time.Since(now))
}

func (t *DynamicSecretsTestSuite) TestReadAndDec() {
	t.InitLocalConfigs()
	pubKey := t.Tc.LocalAgePubkey
	privKey := t.Tc.LocalAgePkey
	age := encryption.NewAge(privKey, pubKey)
	p := filepaths.Path{
		DirIn:  "keygen",
		DirOut: "keygen",
		FnIn:   "key-2.txt.age",
	}
	val, err := ReadAddress(ctx, p, t.S3, age)
	t.Require().NoError(err)
	t.Require().NotEmpty(val)
}

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
	sec, err := LookupAndCreateSecrets(ctx, t.Tc.ProductionLocalTemporalOrgID, topName, kns)
	t.Require().NoError(err)
	t.Require().NotNil(sec)
	//t.Require().NotNil(sec.StringData)
	//t.Assert().Equal(t.Tc.MainnetNodeUrl, sec.StringData["rpc"])
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

	sec, err := CreateSecrets(ctx, zeus_common_types.CloudCtxNs{
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

}

func (t *DynamicSecretsTestSuite) TestPackageSecretProd() {
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

	sec, err := CreateSecrets(ctx, zeus_common_types.CloudCtxNs{
		CloudProvider: "",
		Region:        "",
		Context:       "",
		Namespace:     "test",
		Alias:         "",
		Env:           "",
	}, aegis_secrets.OrgSecretRef{
		OrgSecretReferences: autogen_bases.OrgSecretReferences{
			SecretName: "postgres-auth",
		},
		OrgSecretKeyValReferencesSlice: []autogen_bases.OrgSecretKeyValReferences{
			{
				SecretEnvVarRef: "PG_CONN_STR",
				SecretNameRef:   "postgres-conn-str",
				SecretKeyRef:    auth_startup.PgSecret,
			},
		},
	})

	/*
	   - name: PG_CONN_STR
	     valueFrom:
	       secretKeyRef:
	         name: postgres-auth
	         key: postgres-conn-str
	*/
	t.Require().NoError(err)
	t.Require().NotNil(sec)
}

func TestDynamicSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(DynamicSecretsTestSuite))
}
