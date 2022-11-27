package beacon_cookbook

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
	CloudCtxNs: beaconCloudCtxNs,
}

// chart workload metadata
var execClientChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "geth",
	ChartName:        "geth-hercules",
	ChartDescription: "geth-hercules",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var beaconExecClientChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./ethereum/beacon/infra/exec_client",
	DirOut:      "./ethereum/outputs",
	FnIn:        "geth", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
