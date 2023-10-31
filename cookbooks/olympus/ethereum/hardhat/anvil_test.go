package olympus_hardhat

import (
	"context"
	"fmt"
)

func (t *HardhatCookbookTestSuite) TestCreateClusterClassAnvil() {
	cd := anvilClusterDefinition
	gcd := cd.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	err := gcd.CreateClusterClassDefinitions(context.Background(), t.ZeusTestClient)
	t.Require().Nil(err)
}

func (t *HardhatCookbookTestSuite) TestChartUploadAnvil() {
	cd := anvilClusterDefinition
	_, err := cd.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(err)
}

func (t *HardhatCookbookTestSuite) TestCreateClusterClassAnvilServerless() {
	gcd := serverlessAnvilClusterDefinition.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	err := gcd.CreateClusterClassDefinitions(context.Background(), t.ZeusTestClient)
	t.Require().Nil(err)
}

func (t *HardhatCookbookTestSuite) TestChartUploadAnvilServerless() {
	_, err := serverlessAnvilClusterDefinition.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(err)
}

func (t *HardhatCookbookTestSuite) TestChartUploadAnvilServerlessDev() {
	anvilServerlessChartPath.DirIn = "./olympus/ethereum/hardhat/serverless_anvil_dev"
	anvilServerlessChartPath.FnIn = "anvil-serverless-dev"
	cd := serverlessAnvilClusterDefinition
	cd.ClusterClassName = "anvil-serverless-dev"
	_, err := cd.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(err)
}
