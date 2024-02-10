package zeus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var (
	AgeEnc  = encryption.Age{}
	K8Util  autok8s_core.K8Util
	KeysCfg auth_startup.AuthConfig
)

func VerifyClusterAuthAndGetKubeCfg(ctx context.Context, ou org_users.OrgUser, cloudCtxNs zeus_common_types.CloudCtxNs) (autok8s_core.K8Util, error) {
	p, err := authorized_clusters.SelectAuthedClusterByRouteAndOrgID(ctx, ou, cloudCtxNs)
	if err != nil {
		return autok8s_core.K8Util{}, err
	}
	if p == nil {
		return autok8s_core.K8Util{}, fmt.Errorf("no auth found")
	}
	return GetKubeConfig(ctx, ou, *p)
}

func VerifyClusterAuthAndGetKubeCfgPtr(ctx context.Context, ou org_users.OrgUser, cloudCtxNs zeus_common_types.CloudCtxNs) (*autok8s_core.K8Util, error) {
	p, err := authorized_clusters.SelectAuthedClusterByRouteAndOrgID(ctx, ou, cloudCtxNs)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	k, err := GetKubeConfig(ctx, ou, *p)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("CheckKubeConfig: ConnectToK8sFromInMemFsCfgPathOrErr")
		return nil, err
	}
	return &k, err
}

func GetKubeConfig(ctx context.Context, ou org_users.OrgUser, p authorized_clusters.K8sClusterConfig) (autok8s_core.K8Util, error) {
	k := autok8s_core.K8Util{}
	_, perr := aws_secrets.GetServiceAccountSecrets(ctx, ou)
	if perr != nil {
		log.Err(perr).Interface("ou", ou).Msg("GetExtClusterConfigs: GetServiceAccountSecrets")
		return k, perr
	}
	inKey := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String(p.GetHashedKey(ou.OrgID) + ".kube.tar.gz.age"),
	}
	inMemFsK8sConfig := auth_startup.ExtClustersRunDigitalOceanS3BucketObjAuthProcedure(ctx, KeysCfg, inKey)
	err := k.ConnectToK8sFromInMemFsCfgPathOrErr(inMemFsK8sConfig)
	if err != nil {
		log.Err(err).Interface("key", p.GetHashedKey(ou.OrgID)).Msg("CheckKubeConfig: ConnectToK8sFromInMemFsCfgPathOrErr")
		return autok8s_core.K8Util{}, err
	}
	return k, err
}

func VerifyClusterAuthFromCtxOnlyAndGetKubeCfg(ctx context.Context, ou org_users.OrgUser, cloudCtxNs zeus_common_types.CloudCtxNs) (*autok8s_core.K8Util, error) {
	p, err := authorized_clusters.SelectAuthedClusterByRouteOnlyAndOrgID(ctx, ou, cloudCtxNs)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	k, err := GetKubeConfig(ctx, ou, *p)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("CheckKubeConfig: ConnectToK8sFromInMemFsCfgPathOrErr")
		return nil, err
	}
	return &k, nil
}
