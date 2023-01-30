package hestia_olympus_cookbook

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var HestiaUploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:      "hestia",
	ChartName:         "hestia",
	ChartDescription:  "hestia",
	SkeletonBaseName:  "hestia",
	ComponentBaseName: "hestia",
	ClusterClassName:  "olympus",
	Tag:               "latest",
	Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var HestiaCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "hestia", // set with your own namespace
	Env:           "production",
}

var HestiaDeployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: HestiaCloudCtxNs,
}

var HestiaChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/hestia/infra",
	DirOut:      "./olympus/outputs",
	FnIn:        "hestia", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
