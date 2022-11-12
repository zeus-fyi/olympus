package zeus_client

import (
	base_rest_client "github.com/zeus-fyi/olympus/pkg/iris/resty_base"
)

type ZeusClient struct {
	base_rest_client.Resty
}

func NewZeusClient(baseURL, bearer string) ZeusClient {
	z := ZeusClient{}
	z.Resty = base_rest_client.GetBaseRestyAresTestClient(baseURL, bearer)
	return z
}
