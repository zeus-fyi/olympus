package poseidon_olympus_cookbook

import (
	"fmt"
)

func (t *PoseidonCookbookTestSuite) TestDeploy() {
	_, rerr := PoseidonClusterDefinition.UploadChartsFromClusterDefinition(ctx, t.ZeusExtTestClient, true)
	t.Require().Nil(rerr)
	//cdep := PoseidonClusterDefinition.GenerateDeploymentRequest()

	//_, err := t.ZeusExtTestClient.DeployCluster(ctx, cdep)
	//t.Require().Nil(err)
}

func (t *PoseidonCookbookTestSuite) TestCreateClusterClass() {
	gcd := PoseidonClusterDefinition.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	err := gcd.CreateClusterClassDefinitions(ctx, t.ZeusExtTestClient)
	t.Require().Nil(err)
}
