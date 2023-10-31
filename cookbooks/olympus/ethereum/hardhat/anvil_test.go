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
	cd := serverlessAnvilClusterDefinition
	cd.ClusterClassName = "anvil-serverless-dev"
	_, err := serverlessAnvilClusterDefinition.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(err)
}
