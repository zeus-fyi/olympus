package hestia_cluster_configs

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hestia_eks_aws "github.com/zeus-fyi/olympus/pkg/hestia/aws"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

func GetExtClusterConfigs(ctx context.Context, ou org_users.OrgUser) ([]authorized_clusters.K8sClusterConfig, error) {
	var extClusterConfigs []authorized_clusters.K8sClusterConfig

	ps, perr := aws_secrets.GetServiceAccountSecrets(ctx, ou)
	if perr != nil {
		log.Err(perr).Interface("ou", ou).Msg("GetExtClusterConfigs: GetServiceAccountSecrets")
		return nil, perr
	}
	for clusterName, creds := range ps.AwsEksServiceMap {
		eksCredsAuth := hestia_eks_aws.EksCredentials{
			Creds:       creds,
			ClusterName: clusterName,
		}
		kubeConfig, err := hestia_eks_aws.GetKubeConfig(ctx, eksCredsAuth)
		if err != nil {
			log.Err(err).Interface("ou", ou).Msg("GetExtClusterConfigs: GetKubeConfig")
			return nil, err
		}
		kubeConfigYAML, err := yaml.Marshal(&kubeConfig)
		if err != nil {
			log.Err(err).Interface("ou", ou).Msg("GetExtClusterConfigs: GetKubeConfig")
			return nil, err
		}
		p := filepaths.Path{
			PackageName: "",
			DirIn:       "/.kube",
			FnIn:        "config",
		}
		k := zeus_core.K8Util{}
		inMemFilestore := memfs.NewMemFs()
		err = inMemFilestore.MakeFileIn(&p, kubeConfigYAML)
		if err != nil {
			log.Err(err).Interface("ou", ou).Msg("GetExtClusterConfigs: MakeFileIn")
		}
		k.ConnectToK8sFromInMemFsCfgPath(inMemFilestore)
		ctxes, err := k.GetContexts()
		if err != nil {
			log.Err(err).Interface("ou", ou).Msg("GetExtClusterConfigs: GetContexts")
			return nil, err
		}
		for name, _ := range ctxes {
			ec := authorized_clusters.K8sClusterConfig{
				CloudProvider:     "aws",
				Region:            creds.Region,
				Context:           name,
				ContextAlias:      clusterName,
				Env:               "production",
				IsActive:          false,
				InMemFsKubeConfig: inMemFilestore,
				Path:              p,
			}
			extClusterConfigs = append(extClusterConfigs, ec)
		}

		cmp := compression.NewCompression()
		err = cmp.GzipCompressDir(&p)
	}
	return extClusterConfigs, perr
}
