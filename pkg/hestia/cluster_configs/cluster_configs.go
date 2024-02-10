package hestia_cluster_configs

import (
	"context"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hestia_eks_aws "github.com/zeus-fyi/olympus/pkg/hestia/aws"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
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
			Ou:          ou,
		}
		_, kubeConfig, err := hestia_eks_aws.GetEksKubeConfig(ctx, eksCredsAuth)
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
			zctx := zeus_common_types.CloudCtxNs{
				CloudProvider: "aws",
				Region:        creds.Region,
				Context:       name,
				Namespace:     "kube-system",
				Env:           "production",
			}
			nses, berr := k.GetNamespaces(ctx, zctx)
			if berr != nil {
				log.Err(berr).Interface("nses", nses).Msg("GetExtClusterConfigs: GetNamespaces")
				return nil, berr
			}
			for _, nv := range nses.Items {
				fmt.Println(nv.Name)
			}

			ec := authorized_clusters.K8sClusterConfig{
				CloudCtxNs:        zctx,
				ContextAlias:      clusterName,
				IsActive:          false,
				InMemFsKubeConfig: inMemFilestore,
				Path:              p,
			}

			extClusterConfigs = append(extClusterConfigs, ec)
		}
	}
	return extClusterConfigs, perr
}

type AwsAuthConfigMap struct {
	MapRoles []RoleEntry `json:"mapRoles"`
	MapUsers []UserEntry `json:"mapUsers"`
}

type UserEntry struct {
	UserARN  string   `json:"userarn"`
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
}

type RoleEntry struct {
	Groups   []string `json:"groups"`
	RoleARN  string   `json:"rolearn"`
	Username string   `json:"username"`
}

//cms, cerr := k.GetConfigMapWithKns(ctx, zctx, "aws-auth", nil)
//if cerr != nil {
//	log.Err(cerr).Interface("ou", ou).Msg("GetExtClusterConfigs: GetContexts")
//	return nil, cerr
//}
//var awsAuthMapRoles AwsAuthConfigMap
//if cms.Data != nil {
//	_, ok := cms.Data["mapUsers"]
//	if !ok {
//		eka, eerr := hestia_eks_aws.InitAwsEKS(ctx, eksCredsAuth.Creds)
//		if eerr != nil {
//			log.Err(eerr).Interface("ou", ou).Msg("GetExtClusterConfigs: GetContexts")
//			return nil, eerr
//		}
//		awsAuthMapRoles.MapUsers = []UserEntry{
//			{
//				UserARN:  aws.StringValue(eka.Arn),
//				Username: eka.Username,
//				Groups:   []string{"system:masters"},
//			},
//		}
//		b, berr := yaml.Marshal(awsAuthMapRoles.MapUsers)
//		if berr != nil {
//			log.Err(berr).Interface("ou", ou).Msg("GetExtClusterConfigs: GetContexts")
//			return nil, berr
//		}
//		cms.Data["mapUsers"] = string(b)
//		cms2, kerr := k.UpdateConfigMapWithKns(ctx, zctx, cms, nil)
//		if kerr != nil {
//			log.Err(kerr).Interface("ou", ou).Interface("cms2", cms2).Msg("GetExtClusterConfigs: GetContexts")
//			return nil, kerr
//		}
//	}
//}
