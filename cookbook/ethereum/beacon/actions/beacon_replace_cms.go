package beacon_actions

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_resp_types"
)

func (b *BeaconActionsClient) ReplaceConfigsConsensusClient(ctx context.Context, tar zeus_req_types.TopologyDeployRequest) (zeus_resp_types.TopologyDeployStatus, error) {
	b.ConfigPaths.FnIn = b.ConsensusClient
	b.ConfigPaths.DirIn += "/consensus_client/alt_configs"
	resp, err := b.DeployReplace(ctx, b.ConfigPaths, tar)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("ReplaceConfigsConsensusClient")
		return zeus_resp_types.TopologyDeployStatus{}, err
	}
	return resp, err
}

func (b *BeaconActionsClient) ReplaceConfigsExecClient(ctx context.Context, tar zeus_req_types.TopologyDeployRequest) (zeus_resp_types.TopologyDeployStatus, error) {
	b.ConfigPaths.FnIn = b.ExecClient
	b.ConfigPaths.DirIn += "/exec_client/alt_configs"
	resp, err := b.DeployReplace(ctx, b.ConfigPaths, tar)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("ReplaceConfigsExecClient")
		return zeus_resp_types.TopologyDeployStatus{}, err
	}
	return resp, err
}
