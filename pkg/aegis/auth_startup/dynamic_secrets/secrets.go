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

func LookupAndCreateSecret(ctx context.Context, orgID int, topName string, kns zeus_common_types.CloudCtxNs) (*v1.Secret, error) {
	orgSecrets, err := aegis_secrets.SelectOrgSecretRef(ctx, orgID, topName)
	if err != nil {
		return nil, err
	}
	return CreateSecret(ctx, kns, orgSecrets)
}

func CreateSecret(ctx context.Context, kns zeus_common_types.CloudCtxNs, os aegis_secrets.OrgSecretRef) (*v1.Secret, error) {
	var sec *v1.Secret
	for _, kv := range os.OrgSecretKeyValReferencesSlice {
		value, err := getSecretValue(ctx, kv.SecretKeyRef)
		if err != nil {
			return nil, err
		}
		sec = zeus_core.CreateSecretWrapper(sec, kns, os.SecretName, kv.SecretNameRef, value)
	}
	return sec, nil
}

func getSecretValue(ctx context.Context, refName string) (string, error) {
	secrets := auth_startup.SecretsWrapper{}
	return secrets.ReadSecret(ctx, AegisInMemSecrets, refName)
}
