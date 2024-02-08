package orchestrate_pods_activities

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type PodsActivity struct {
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (p *PodsActivity) GetActivities() ActivitiesSlice {
	return []interface{}{p.DeletePod}
}

func (p *PodsActivity) DeletePod(ctx context.Context, ou org_users.OrgUser, podName string, cctx zeus_common_types.CloudCtxNs) error {
	k, err := zeus.VerifyClusterAuthFromCtxOnlyAndGetKubeCfg(ctx, ou, cctx)
	if err != nil {
		log.Warn().Interface("ou", ou).Interface("req", cctx).Msg("PodsCloudCtxNsMiddleware: IsOrgCloudCtxNsAuthorizedFromID")
		return err
	}
	kut := zeus.K8Util
	if k != nil {
		kut = *k
	}
	err = kut.DeleteFirstPodLike(ctx, cctx, podName, nil, nil)
	if err != nil {
		log.Err(err).Msg("PodsActivity: DeletePod: DeleteFirstPodLike")
		return err
	}
	return nil
}
