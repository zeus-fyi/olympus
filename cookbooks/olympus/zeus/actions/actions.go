package zeus_actions

import (
	zeus_cookbook "github.com/zeus-fyi/olympus/cookbooks/olympus/zeus"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
)

type ZeusActionsClient struct {
	zeus_client.ZeusClient
}

var zeusBasePar = zeus_pods_reqs.PodActionRequest{
	TopologyDeployRequest: zeus_cookbook.ZeusDeployKnsReq,
	PodName:               "zeus",
	FilterOpts:            nil,
	ClientReq:             nil,
	DeleteOpts:            nil,
}
