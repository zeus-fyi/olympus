package apollo_client

import (
	base_rest_client "github.com/zeus-fyi/olympus/pkg/iris/resty_base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
)

type Apollo struct {
	zeus_client.ZeusClient
}

func NewApollo(baseURL, bearer string) Apollo {
	z := Apollo{}
	z.Resty = base_rest_client.GetBaseRestyClient(baseURL, bearer)
	return z
}

const ApolloEndpoint = "https://apollo.eth.zeus.fyi"

func NewDefaultApolloClient(bearer string) Apollo {
	return NewApollo(ApolloEndpoint, bearer)
}

const ApolloLocalEndpoint = "http://localhost:9000"

func NewLocalApolloClient(bearer string) Apollo {
	return NewApollo(ApolloLocalEndpoint, bearer)
}
