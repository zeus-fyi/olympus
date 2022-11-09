package test

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/pretty"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/deploy/workload_state"
)

var updateStatusHost = "http://localhost:9001/v1/internal/deploy/status"

var topologyID = 1667886918118352000

func UpdateDeploymentStatusApiCall() error {
	cfg := configs.InitLocalTestConfigs()

	wsr := workload_state.InternalWorkloadStatusUpdateRequest{
		TopologyID:     topologyID,
		TopologyStatus: "InProgress",
	}
	topologyActionRequestPayload, err := json.Marshal(wsr)
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
		SetBody(wsr).
		Post(updateStatusHost)

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
