package dynamic_secrets

import (
	"context"

	aegis_secrets "github.com/zeus-fyi/olympus/datastores/postgres/apps/aegis"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
)

var (
	AegisInMemSecrets memfs.MemFS
)

func LookupAndCreateSecrets(ctx context.Context, orgID int, topName string, kns zeus_common_types.CloudCtxNs) ([]*v1.Secret, error) {
	orgSecrets, err := aegis_secrets.SelectOrgSecretRef(ctx, orgID, topName)
	if err != nil {
		return nil, err
	}
	return CreateSecrets(ctx, kns, orgSecrets)
}

func CreateSecrets(ctx context.Context, kns zeus_common_types.CloudCtxNs, os aegis_secrets.OrgSecretRef) ([]*v1.Secret, error) {
	var sec *v1.Secret
	secrets := make([]*v1.Secret, len(os.OrgSecretKeyValReferencesSlice))
	for i, kv := range os.OrgSecretKeyValReferencesSlice {
		value, err := getSecretValue(ctx, kv.SecretKeyRef)
		if err != nil {
			return nil, err
		}
		sec = zeus_core.CreateSecretWrapper(sec, kns, os.SecretName, kv.SecretNameRef, value)
		secrets[i] = sec
	}
	return secrets, nil
}

func getSecretValue(ctx context.Context, refName string) (string, error) {
	secrets := auth_startup.SecretsWrapper{}
	return secrets.ReadSecret(ctx, AegisInMemSecrets, refName)
}
