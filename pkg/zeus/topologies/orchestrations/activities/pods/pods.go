package orchestrate_pods_activities

import (
	"context"

	"github.com/rs/zerolog/log"
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

func (p *PodsActivity) DeletePod(ctx context.Context, podName string, cctx zeus_common_types.CloudCtxNs) error {
	err := zeus.K8Util.DeleteFirstPodLike(ctx, cctx, podName, nil, nil)
	if err != nil {
		log.Err(err).Msg("DeletePod: DeleteFirstPodLike")
		return err
	}
	return nil
}
