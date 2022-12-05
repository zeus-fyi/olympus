package beacon_cookbooks

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

// set your own topologyID here after uploading a chart workload
var deployExecClientKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: BeaconCloudCtxNs,
}

// chart workload metadata
var ExecClientChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "geth",
	ChartName:        "geth-hercules",
	ChartDescription: "geth-hercules",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),

	SkeletonBaseID: GethSkeletonBaseID,
}

var BeaconExecClientChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./ethereum/beacon/infra/exec_client",
	DirOut:      "./ethereum/outputs",
	FnIn:        "geth", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
