package test

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/pretty"
	"github.com/zeus-fyi/olympus/configs"
)

var readTopMetaHost = "http://localhost:9001/v1/infra/read/topologies"

func ReadTopologiesMetadataAPICall() error {
	cfg := configs.InitLocalTestConfigs()
	client := resty.New()
	resp, err := client.R().
		SetAuthToken(cfg.LocalBearerToken).
		Get(readTopMetaHost)

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
