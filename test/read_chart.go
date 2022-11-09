package test

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/pretty"
	"github.com/zeus-fyi/olympus/configs"
	read_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/read"
)

var readChartHost = "http://localhost:9001/v1/infra/read/chart"

func ReadChartAPICall() error {
	cfg := configs.InitLocalTestConfigs()
	tar := read_infra.TopologyReadRequest{
		TopologyID: 1667887316382350000,
	}
	topologyActionRequestPayload, err := json.Marshal(tar)
	if err != nil {
		return err
	}
	fmt.Println("action request json")
	requestJSON := pretty.Pretty(topologyActionRequestPayload)
	requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
	fmt.Println(string(requestJSON))

	client := resty.New()
	resp, err := client.R().
		SetAuthToken(cfg.LocalBearerToken).
		SetBody(tar).
		Post(readChartHost)

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
