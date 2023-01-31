package olympus_hydra_cookbooks

func (t *HydraCookbookTestSuite) TestHydraSecretsCopy() {
	//s1 := "spaces-auth"
	//s2 := "spaces-key"
	//s3 := "age-auth"
	///*
	//	for mainnet
	//	cd.CloudCtxNs.Namespace = mainnetNamespace
	//	cd.ClusterClassName = "hydraMainnet"
	//*/
	//hydraCtxNsTop := kns.TopologyKubeCtxNs{
	//	TopologyID: 0,
	//	CloudCtxNs: ValidatorCloudCtxNs,
	//}
	//req := internal_reqs.InternalSecretsCopyFromTo{
	//	SecretNames: []string{s1, s2, s3},
	//	FromKns: kns.TopologyKubeCtxNs{
	//		TopologyID: 0,
	//		CloudCtxNs: zeus_common_types.CloudCtxNs{
	//			CloudProvider: "do",
	//			Region:        "sfo3",
	//			Context:       "do-sfo3-dev-do-sfo3-zeus",
	//			Namespace:     "zeus",
	//			Env:           "dev",
	//		},
	//	},
	//	ToKns: hydraCtxNsTop,
	//}
	//
	//err := t.ZeusTestClient.CopySecretsFromToNamespace(context.Background(), req)
	//t.Require().Nil(err)
}
