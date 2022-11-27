package artemis_cookbook

import (
	"context"
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

var ctx = context.Background()

// chart workload metadata
var uploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "artemis",
	ChartName:        "artemis",
	ChartDescription: "artemis",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var ZeusCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "artemis", // set with your own namespace
	Env:           "production",
}

var ArtemisDeployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: ZeusCloudCtxNs,
}

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var artemisChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/artemis/infra",
	DirOut:      "./olympus/outputs",
	FnIn:        "artemis", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
