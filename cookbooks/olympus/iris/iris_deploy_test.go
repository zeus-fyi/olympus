package iris_olympus_cookbook

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func (t *IrisCookbookTestSuite) TestDeploy() {
	_, rerr := irisClusterDefinition.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(rerr)
	//cdep := irisClusterDefinition.GenerateDeploymentRequest()

	//_, err := t.ZeusTestClient.DeployCluster(ctx, cdep)
	//t.Require().Nil(err)
}

func (t *IrisCookbookTestSuite) TestClusterDestroy() {
	d := zeus_req_types.TopologyDeployRequest{
		CloudCtxNs: irisCloudCtxNs,
	}
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, d)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *IrisCookbookTestSuite) TestCreateClusterClass() {
	cd := irisClusterDefinition
	gcd := cd.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	err := gcd.CreateClusterClassDefinitions(context.Background(), t.ZeusTestClient)
	t.Require().Nil(err)
}
