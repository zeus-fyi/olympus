package beacon_cookbooks

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

// set your own topologyID here after uploading a chart workload
var DeployExecClientKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: BeaconCloudCtxNs,
}

var ExecClientChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "gethHercules",
	ChartName:        "gethHercules",
	ChartDescription: "gethHercules",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),

	SkeletonBaseName: "gethHercules",
}

var BeaconExecClientChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./ethereum/beacon/infra/exec_client",
	DirOut:      "./ethereum/outputs",
	FnIn:        "gethHercules", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
