package poseidon_olympus_cookbook

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/internal_reqs"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func (t *PoseidonCookbookTestSuite) TestPoseidonSecretsCopy() {
	s1 := "spaces-auth"
	s2 := "spaces-key"
	s3 := "age-auth"
	req := internal_reqs.InternalSecretsCopyFromTo{
		SecretNames: []string{s1, s2, s3},
		FromKns: kns.TopologyKubeCtxNs{
			TopologyID: 0,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "do",
				Region:        "sfo3",
				Context:       "do-sfo3-dev-do-sfo3-zeus",
				Namespace:     "zeus",
				Env:           "dev",
			},
		},
		ToKns: kns.TopologyKubeCtxNs{
			TopologyID: 0,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "do",
				Region:        "sfo3",
				Context:       "do-sfo3-dev-do-sfo3-zeus",
				Namespace:     "poseidon",
				Env:           "production",
			},
		},
	}
	err := t.ZeusTestClient.CopySecretsFromToNamespace(ctx, req)
	t.Require().Nil(err)
}
