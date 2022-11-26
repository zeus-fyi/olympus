package aegis_olympus_cookbook

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

// chart workload metadata
var uploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "aegis",
	ChartName:        "aegis",
	ChartDescription: "aegis",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var topCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "aegis", // set with your own namespace
	Env:           "production",
}

// set your own topologyID here after uploading a chart workload
var deployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1669423376281749000,
	CloudCtxNs: topCloudCtxNs,
}

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var aegisChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/aegis/infra",
	DirOut:      "./olympus/outputs",
	FnIn:        "aegis", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
