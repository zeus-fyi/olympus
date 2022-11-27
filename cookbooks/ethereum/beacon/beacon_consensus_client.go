package beacon_cookbook

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

// set your own topologyID here after uploading a chart workload
var deployConsensusClientKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: beaconCloudCtxNs,
}

// chart workload metadata
var consensusClientChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "lighthouse",
	ChartName:        "lighthouse-hercules",
	ChartDescription: "lighthouse-hercules",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var beaconConsensusClientChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./ethereum/beacon/infra",
	DirOut:      "./ethereum/outputs",
	FnIn:        "lighthouse", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
