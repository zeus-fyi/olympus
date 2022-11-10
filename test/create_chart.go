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

var createChartHost = "https://api.zeus.fyi/v1/infra/create"

var p = structs.Path{
	PackageName: "",
	DirIn:       "../configs/apps/demo",
	DirOut:      "../test",
	Fn:          "demo",
	FnOut:       "demo.tar.gz",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}

func CreateChartApiCall() error {
	cfg := configs.InitLocalTestConfigs()
	var ts chronos.Chronos
	tar := create_infra.TopologyCreateRequest{
		TopologyName:     "demo",
		ChartName:        "demo",
		ChartDescription: "zeus infra demo test",
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

	err = ZipK8sChart(p)
	if err != nil {
		return err
	}

	fp := p.V2FileOutPath()
	fmt.Println("filepath")
	fmt.Println(fp)

	client := resty.New()
	resp, err := client.R().
		SetAuthToken(cfg.ProductionLocalBearerToken).
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

func ZipK8sChart(p structs.Path) error {
	comp := compression.NewCompression()
	err := comp.CreateTarGzipArchiveDir(&p)
	return err
}
