package test

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/pretty"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	topology_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
	create_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

var host = "http://localhost:9001/v1/infra"

var kns = autok8s_core.KubeCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	CtxType:       "dev-sfo3-zeus",
	Namespace:     "demo",
}

var p = structs.Path{
	PackageName: "",
	DirIn:       "../configs/apps/zeus",
	DirOut:      "../test",
	Fn:          "zeus",
	FnOut:       "zeus.tar.gz",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}

func CallAPI() error {

	cfg := configs.InitLocalTestConfigs()
	var ts chronos.Chronos
	tar := create_infra.TopologyActionCreateRequest{
		TopologyActionRequest: base.TopologyActionRequest{
			Action: "create",
			TopologyActivityRequest: topology_activities.TopologyActivityRequest{
				TopologyID: 0,
				Bearer:     cfg.LocalBearerToken,
				Kns:        kns,
				Host:       host,
				NativeK8s:  chart_workload.NativeK8s{},
			},
		},
		TopologyCreateRequest: create_infra.TopologyCreateRequest{
			TopologyName:     "zeus",
			ChartName:        "zeus",
			ChartDescription: "zeus infra dev test",
			Version:          fmt.Sprintf("v0.0.%d", ts.UnixTimeStampNow()),
		},
	}

	topologyActionRequestPayload, err := json.Marshal(tar)
	if err != nil {
		return err
	}
	fmt.Println("action request json")
	requestJSON := pretty.Pretty(topologyActionRequestPayload)
	requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
	fmt.Println(string(requestJSON))

	fp := p.V2FileOutPath()
	client := resty.New()
	resp, err := client.R().
		SetFile("chart", fp).
		Post(host)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("response json")
	respJSON := pretty.Pretty(resp.Body())
	respJSON = pretty.Color(respJSON, pretty.TerminalStyle)
	fmt.Println(string(respJSON))

	return err
}

func ZipK8sChart() error {
	comp := compression.NewCompression()
	err := comp.CreateTarGzipArchiveDir(&p)
	return err
}
