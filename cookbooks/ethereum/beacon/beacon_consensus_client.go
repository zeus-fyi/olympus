package beacon_cookbooks

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

// DeployConsensusClientKnsReq set your own topologyID here after uploading a chart workload
var DeployConsensusClientKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: BeaconCloudCtxNs,
}

// ConsensusClientChart chart workload metadata
var ConsensusClientChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "lighthouse",
	ChartName:        "lighthouse-hercules",
	ChartDescription: "lighthouse-hercules",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),

	SkeletonBaseName: "lighthouse",
}

// BeaconConsensusClientChartPath DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var BeaconConsensusClientChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./ethereum/beacon/infra",
	DirOut:      "./ethereum/outputs",
	FnIn:        "lighthouse", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
