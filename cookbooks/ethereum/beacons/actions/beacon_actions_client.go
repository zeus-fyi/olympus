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
	PrintPath       filepaths.Path
	ConfigPaths     filepaths.Path
	ConsensusClient string
	ExecClient      string
}

// set your own topologyID here after uploading a chart workload
var beaconKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1669159384971627008,
	CloudCtxNs: beaconCloudCtxNs,
}

var beaconCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "beacon", // set with your own namespace
	Env:           "dev",
}

var basePar = zeus_pods_reqs.PodActionRequest{
	TopologyDeployRequest: beaconKnsReq,
	PodName:               "",
	FilterOpts:            nil,
	ClientReq:             nil,
	DeleteOpts:            nil,
}

func NewBeaconActionsClient(baseURL, bearer string) BeaconActionsClient {
	z := BeaconActionsClient{}
	z.Resty = base_rest_client.GetBaseRestyClient(baseURL, bearer)
	return z
}

const ZeusEndpoint = "https://api.zeus.fyi"

func NewDefaultBeaconActionsClient(bearer string) BeaconActionsClient {
	return NewBeaconActionsClient(ZeusEndpoint, bearer)
}

const ZeusLocalEndpoint = "http://localhost:9000"

func NewLocalBeaconActionsClient(bearer string) BeaconActionsClient {
	return NewBeaconActionsClient(ZeusLocalEndpoint, bearer)
}
