package destroy_deploy_activities

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	topology_auths "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/auth"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
)

func (d *DestroyDeployTopologyActivities) DestroyNamespace(ctx context.Context, params base_request.InternalDeploymentActionRequest) error {
	k := topology_auths.K8Util
	p, perr := authorized_clusters.SelectAuthedClusterByRouteOnlyAndOrgID(ctx, params.OrgUser, params.Kns.CloudCtxNs)
	if perr != nil {
		return perr
	}
	if p != nil && !p.IsPublic {
		kp, err := GetKubeConfig(ctx, params.OrgUser, *p)
		if err != nil {
			log.Err(err).Interface("params", params).Msg("DestroyNamespace: CheckKubeConfig: ConnectToK8sFromInMemFsCfgPathOrErr")
			return err
		}
		k = kp
	}
	err := k.DeleteNamespace(ctx, params.Kns.CloudCtxNs)
	if err != nil {
		log.Err(err).Interface("params", params).Msg("DestroyDeployTopologyActivities: DestroyNamespace")
		return err
	}
	return nil
}

func GetKubeConfig(ctx context.Context, ou org_users.OrgUser, p authorized_clusters.K8sClusterConfig) (autok8s_core.K8Util, error) {
	inKey := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String(p.GetHashedKey(ou.OrgID) + ".kube.tar.gz.age"),
	}
	inMemFsK8sConfig := auth_startup.ExtClustersRunDigitalOceanS3BucketObjAuthProcedure(ctx, topology_auths.KeysCfg, inKey)
	k := autok8s_core.K8Util{}
	err := k.ConnectToK8sFromInMemFsCfgPathOrErr(inMemFsK8sConfig)
	if err != nil {
		log.Err(err).Interface("key", p.GetHashedKey(ou.OrgID)).Msg("CheckKubeConfig: ConnectToK8sFromInMemFsCfgPathOrErr")
		return autok8s_core.K8Util{}, err
	}
	return k, err
}
