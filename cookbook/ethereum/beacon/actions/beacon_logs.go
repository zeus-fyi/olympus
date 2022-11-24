package beacon_actions

import (
	"context"
	"path"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
	v1 "k8s.io/api/core/v1"
)

func (b *BeaconActionsClient) PrintConsensusClientPodLogs(ctx context.Context, par zeus_pods_reqs.PodActionRequest) ([]byte, error) {
	b.PrintReqJson(par)
	par.PodName = b.ConsensusClient
	par.ContainerName = b.ConsensusClient
	filter := string_utils.FilterOpts{Contains: b.ConsensusClient}
	logOpts := &v1.PodLogOptions{Container: b.ConsensusClient}
	par.LogOpts = logOpts
	par.FilterOpts = &filter

	resp, err := b.GetPodLogs(ctx, par)
	if err != nil {
		return nil, err
	}
	b.PrintPath.FnOut = b.ConsensusClient + "_logs"
	b.PrintPath.DirOut = path.Join(b.PrintPath.DirIn, "/consensus_client")
	err = b.PrintPath.Print(resp)
	return resp, err
}

func (b *BeaconActionsClient) PrintExecClientPodLogs(ctx context.Context, par zeus_pods_reqs.PodActionRequest) ([]byte, error) {
	b.PrintReqJson(par)
	par.PodName = b.ExecClient
	par.ContainerName = b.ExecClient
	logOpts := &v1.PodLogOptions{Container: b.ExecClient}
	par.LogOpts = logOpts

	filter := string_utils.FilterOpts{Contains: b.ConsensusClient}
	par.FilterOpts = &filter
	resp, err := b.GetPodLogs(ctx, par)
	if err != nil {
		return nil, err
	}
	b.PrintPath.FnOut = b.ExecClient + "_logs"
	b.PrintPath.DirOut = path.Join(b.PrintPath.DirIn, "/exec_client")
	err = b.PrintPath.WriteToFileOutPath(resp)
	return resp, err
}
