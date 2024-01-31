package aegis_olympus_cookbook

import (
	"fmt"
)

const (
	aegis   = "aegis"
	olympus = "olympus"
)

func (t *AegisCookbookTestSuite) TestDeploy() {
	_, rerr := AegisClusterDefinition.UploadChartsFromClusterDefinition(ctx, t.ZeusExtTestClient, true)
	t.Require().Nil(rerr)
}

func (t *AegisCookbookTestSuite) TestCreateClusterClass() {
	gcd := AegisClusterDefinition.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	err := gcd.CreateClusterClassDefinitions(ctx, t.ZeusExtTestClient)
	t.Require().Nil(err)
}
