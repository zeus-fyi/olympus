package aegis_olympus_cookbook

import (
	"context"

	olympus_beacon_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/beacons"
	olympus_hydra_cookbooks "github.com/zeus-fyi/olympus/cookbooks/olympus/ethereum/hydra"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/internal_reqs"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

func (t *AegisCookbookTestSuite) TestAegisSecretsCopy() {
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
				Namespace:     "aegis",
				Env:           "dev",
			},
		},
	}
	err := t.ZeusTestClient.CopySecretsFromToNamespace(ctx, req)
	t.Require().Nil(err)
}

func (t *AegisCookbookTestSuite) TestAegisSecretsCopyTxFetcher() {
	s1 := "postgres-auth"
	s2 := "dynamodb-auth"
	req := internal_reqs.InternalSecretsCopyFromTo{
		SecretNames: []string{s1, s2},
		FromKns: kns.TopologyKubeCtxNs{
			TopologyID: 0,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "do",
				Region:        "sfo3",
				Context:       "do-sfo3-dev-do-sfo3-zeus",
				Namespace:     "ethereum",
				Env:           "dev",
			},
		},
		ToKns: kns.TopologyKubeCtxNs{
			TopologyID: 0,
			CloudCtxNs: zeus_common_types.CloudCtxNs{
				CloudProvider: "do",
				Region:        "sfo3",
				Context:       "do-sfo3-dev-do-sfo3-zeus",
				Namespace:     "tx-fetcher",
				Env:           "dev",
			},
		},
	}
	err := t.ZeusTestClient.CopySecretsFromToNamespace(ctx, req)
	t.Require().Nil(err)
}

func (t *AegisCookbookTestSuite) TestAegisSecretsCopyPG() {
	s1 := "postgres-auth"
	req := internal_reqs.InternalSecretsCopyFromTo{
		SecretNames: []string{s1},
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
				Namespace:     "ethereum",
				Env:           "dev",
			},
		},
	}
	err := t.ZeusTestClient.CopySecretsFromToNamespace(ctx, req)
	t.Require().Nil(err)
}

func (t *AegisCookbookTestSuite) TestMainnetBeaconSecretsCopy() {
	s1 := "spaces-auth"
	s2 := "spaces-key"
	s3 := "age-auth"
	mainnetBeaconCtxNsTop := kns.TopologyKubeCtxNs{
		TopologyID: 0,
		CloudCtxNs: olympus_beacon_cookbooks.GetBeaconCloudCtxNs(hestia_req_types.Mainnet),
	}

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
		ToKns: mainnetBeaconCtxNsTop,
	}

	err := t.ZeusTestClient.CopySecretsFromToNamespace(context.Background(), req)
	t.Require().Nil(err)
}

func (t *AegisCookbookTestSuite) TestHydraSecretsCopy() {
	s1 := "spaces-auth"
	s2 := "spaces-key"
	s3 := "age-auth"
	/*
		for mainnet
		cd.CloudCtxNs.Namespace = mainnetNamespace
		cd.ClusterClassName = "hydraMainnet"
	*/

	hydraCtxNsTop := kns.TopologyKubeCtxNs{
		TopologyID: 0,
		CloudCtxNs: olympus_hydra_cookbooks.ValidatorCloudCtxNs,
	}
	/*
		for goerli
		hydraCtxNsTop.Namespace = "goerli-staking"
		cd.ClusterClassName = "hydraGoerli"
	*/

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
		ToKns: hydraCtxNsTop,
	}

	err := t.ZeusTestClient.CopySecretsFromToNamespace(context.Background(), req)
	t.Require().Nil(err)
}
