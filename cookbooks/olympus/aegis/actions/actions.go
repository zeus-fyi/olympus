package aegis_actions

import (
	aegis_olympus_cookbook "github.com/zeus-fyi/olympus/cookbooks/olympus/aegis"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
)

type AegisActionsClient struct {
	zeus_client.ZeusClient
}

var aegisBasePar = zeus_pods_reqs.PodActionRequest{
	TopologyDeployRequest: aegis_olympus_cookbook.AegisDeployKnsReq,
	PodName:               "aegis",
	FilterOpts:            nil,
	ClientReq:             nil,
	DeleteOpts:            nil,
}
