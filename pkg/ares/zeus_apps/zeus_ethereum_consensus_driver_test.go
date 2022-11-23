package zeus_apps

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

// set your own topologyID here after uploading a chart workload
var deployConsensusKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1668825071425428000,
	CloudCtxNs: topCloudCtxNs,
}

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var consensusChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./mocks/kubernetes_apps/beacon/consensus_client",
	DirOut:      "./mocks/kubernetes_apps/beacon/consensus_client_out",
	FnIn:        "consensus_client", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
