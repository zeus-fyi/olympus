package test

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/pretty"
	"github.com/zeus-fyi/olympus/configs"
	destroy_deploy_request "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/destroy"
)

var destroyDeployDemoChartHostProduction = "https://api.zeus.fyi/v1/deploy/destroy"

func DestroyDeployDemoProdChartApiCall() error {
	cfg := configs.InitLocalTestConfigs()

	deployKns := destroy_deploy_request.TopologyDestroyDeployRequest{
		TopologyID:    1667958167340986000,
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "dev-sfo3-zeus",
		Namespace:     "demo",
		Env:           "dev",
	}

	topologyActionRequestPayload, err := json.Marshal(deployKns)
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
		SetBody(deployKns).
		Post(destroyDeployDemoChartHostProduction)

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
