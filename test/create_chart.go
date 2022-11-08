package test

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/pretty"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	create_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"
)

var createChartHost = "http://localhost:9001/v1/infra/create"

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
	tar := create_infra.TopologyCreateRequest{
		TopologyName:     "zeus",
		ChartName:        "zeus",
		ChartDescription: "zeus infra dev test",
		Version:          fmt.Sprintf("v0.0.%d", ts.UnixTimeStampNow()),
	}

	topologyActionRequestPayload, err := json.Marshal(tar)
	if err != nil {
		return err
	}
	fmt.Println("action request json")
	requestJSON := pretty.Pretty(topologyActionRequestPayload)
	requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
	fmt.Println(string(requestJSON))

	err = ZipK8sChart()
	if err != nil {
		return err
	}

	fp := p.V2FileOutPath()
	fmt.Println("filepath")
	fmt.Println(fp)

	client := resty.New()
	resp, err := client.R().
		SetAuthToken(cfg.LocalBearerToken).
		SetFormData(map[string]string{
			"topologyName":     tar.TopologyName,
			"chartName":        tar.ChartName,
			"chartDescription": tar.ChartDescription,
			"version":          tar.Version,
		}).
		SetFile("chart", fp).
		Post(createChartHost)

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
