package zeus_client

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/internal_reqs"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func (t *ZeusClientTestSuite) TestSecretsCopy() {
	s1 := "dynamodb-auth"
	req := internal_reqs.InternalSecretsCopyFromTo{
		SecretNames: []string{s1},
		FromKns: kns.TopologyKubeCtxNs{
			TopologyID: 0,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "ovh",
				Region:        "us-west-or-1",
				Context:       "kubernetes-admin@zeusfyi",
				Namespace:     "c9bbfe9e-9922-4d25-bc31-e7c776cfb349",
				Env:           "dev",
			},
		},
		ToKns: kns.TopologyKubeCtxNs{
			TopologyID: 0,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "ovh",
				Region:        "us-west-or-1",
				Context:       "kubernetes-admin@zeusfyi",
				Namespace:     "672973a4-087c-40b0-a9aa-cf2e183cd6c3",
				Env:           "dev",
			},
		},
	}
	err := t.ZeusTestClient.CopySecretsFromToNamespace(ctx, req)
	t.Require().Nil(err)
}
