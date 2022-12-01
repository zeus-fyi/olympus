package zeus_client

import (
	base_rest_client "github.com/zeus-fyi/olympus/pkg/iris/resty_base"
)

type ZeusClient struct {
	base_rest_client.Resty
}

func NewZeusClient(baseURL, bearer string) ZeusClient {
	z := ZeusClient{}
	z.Resty = base_rest_client.GetBaseRestyClient(baseURL, bearer)
	return z
}

const ZeusEndpoint = "https://api.zeus.fyi"

func NewDefaultZeusClient(bearer string) ZeusClient {
	return NewZeusClient(ZeusEndpoint, bearer)
}

const ZeusLocalEndpoint = "http://localhost:9001"

func NewLocalZeusClient(bearer string) ZeusClient {
	return NewZeusClient(ZeusLocalEndpoint, bearer)
}
