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
