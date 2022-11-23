package athena_client

import base_rest_client "github.com/zeus-fyi/olympus/pkg/iris/resty_base"

type AthenaClient struct {
	base_rest_client.Resty
}

func NewAthenaClient(baseURL, bearer string) AthenaClient {
	z := AthenaClient{}
	z.Resty = base_rest_client.GetBaseRestyAresTestClient(baseURL, bearer)
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
