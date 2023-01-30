package hera_olympus_cookbook

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var HeraUploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:      "Hera",
	ChartName:         "Hera",
	ChartDescription:  "Hera",
	SkeletonBaseName:  "hera",
	ComponentBaseName: "hera",
	ClusterClassName:  "olympus",
	Tag:               "latest",
	Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var HeraCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "hera", // set with your own namespace
	Env:           "production",
}

var HeraDeployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: HeraCloudCtxNs,
}

var HeraChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/hera/infra",
	DirOut:      "./olympus/outputs",
	FnIn:        "hera", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
