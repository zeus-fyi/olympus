package poseidon_olympus_cookbook

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

var PoseidonUploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:     "poseidon",
	ChartName:        "poseidon",
	ChartDescription: "poseidon",
	Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var PoseidonCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "poseidon", // set with your own namespace
	Env:           "production",
}

var PoseidonDeployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1670807021640692000,
	CloudCtxNs: PoseidonCloudCtxNs,
}

var PoseidonChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/poseidon/infra",
	DirOut:      "./olympus/outputs",
	FnIn:        "poseidon", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
