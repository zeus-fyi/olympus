package athena_client

import (
	zeus_client "github.com/zeus-fyi/zeus/zeus/z_client"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

type AthenaClient struct {
	zeus_client.ZeusClient
}

func NewAthenaClient(baseURL, bearer string) AthenaClient {
	z := AthenaClient{}
	z.Resty = resty_base.GetBaseRestyClient(baseURL, bearer)
	return z
}

const ZeusEndpoint = "https://api.zeus.fyi"

func NewDefaultAthenaClient(bearer string) AthenaClient {
	return NewAthenaClient(ZeusEndpoint, bearer)
}

const ZeusLocalEndpoint = "http://localhost:9003"

func NewLocalAthenaClient(bearer string) AthenaClient {
	return NewAthenaClient(ZeusLocalEndpoint, bearer)
}
