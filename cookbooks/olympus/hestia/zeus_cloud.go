package hestia_olympus_cookbook

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var ZeusCloudUploadChart = zeus_req_types.TopologyCreateRequest{
	TopologyName:      "zeusCloud",
	ChartName:         "zeusCloud",
	ChartDescription:  "zeusCloud",
	SkeletonBaseName:  "zeusCloud",
	ComponentBaseName: "zeusCloud",
	ClusterClassName:  "olympus",
	Tag:               "latest",
	Version:           fmt.Sprintf("v0.0.%d", time.Now().Unix()),
}

var ZeusCloudCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "hestia", // set with your own namespace
	Env:           "production",
}

var ZeusCloudDeployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 0,
	CloudCtxNs: ZeusCloudCloudCtxNs,
}

var ZeusCloudChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/hestia/frontend_infra",
	DirOut:      "./olympus/outputs",
	FnIn:        "zeusCloud", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
