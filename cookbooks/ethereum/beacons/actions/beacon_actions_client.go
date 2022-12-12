package beacon_actions

import (
	base_rest_client "github.com/zeus-fyi/olympus/pkg/iris/resty_base"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

type BeaconActionsClient struct {
	zeus_client.ZeusClient
	BeaconKnsReq    zeus_req_types.TopologyDeployRequest
	PrintPath       filepaths.Path
	ConfigPaths     filepaths.Path
	ConsensusClient string
	ExecClient      string
}

// BeaconKnsReq set your own topologyID here after uploading a chart workload
var BeaconKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1669159384971627008,
	CloudCtxNs: BeaconCloudCtxNs,
}

var BeaconCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "ethereum", // set with your own namespace
	Env:           "production",
}

var basePar = zeus_pods_reqs.PodActionRequest{
	TopologyDeployRequest: BeaconKnsReq,
	PodName:               "",
	FilterOpts:            nil,
	ClientReq:             nil,
	DeleteOpts:            nil,
}

func NewBeaconActionsClient(baseURL, bearer string, kCtxNs zeus_req_types.TopologyDeployRequest) BeaconActionsClient {
	z := BeaconActionsClient{}
	z.BeaconKnsReq = kCtxNs
	z.Resty = base_rest_client.GetBaseRestyClient(baseURL, bearer)
	return z
}

const ZeusEndpoint = "https://api.zeus.fyi"

func NewDefaultBeaconActionsClient(bearer string, kCtxNs zeus_req_types.TopologyDeployRequest) BeaconActionsClient {
	return NewBeaconActionsClient(ZeusEndpoint, bearer, kCtxNs)
}

const ZeusLocalEndpoint = "http://localhost:9001"

func NewLocalBeaconActionsClient(bearer string) BeaconActionsClient {
	return NewBeaconActionsClient(ZeusLocalEndpoint, bearer, BeaconKnsReq)
}
