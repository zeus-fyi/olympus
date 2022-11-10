package test

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/pretty"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/base_request"
)

var createProdInternalNsHostProduction = "https://api.zeus.fyi/v1/internal/deploy/namespace"

func CreateInternalProdNs() error {
	cfg := configs.InitLocalTestConfigs()

	newKns := kns.TopologyKubeCtxNs{
		TopologiesKns: autogen_bases.TopologiesKns{
			TopologyID:    1667958167340986000,
			CloudProvider: "do",
			Region:        "sfo3",
			Context:       "dev-sfo3-zeus",
			Namespace:     "demo",
			Env:           "dev",
		}}
	req := base_request.InternalDeploymentActionRequest{
		Kns:       newKns,
		OrgUser:   org_users.OrgUser{},
		NativeK8s: chart_workload.NativeK8s{},
	}

	topologyActionRequestPayload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	fmt.Println("action request json")
	requestJSON := pretty.Pretty(topologyActionRequestPayload)
	requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
	fmt.Println(string(requestJSON))

	client := resty.New()
	resp, err := client.R().
		SetAuthToken(cfg.ProductionLocalTemporalBearerToken).
		SetBody(req).
		Post(createProdInternalNsHostProduction)

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
