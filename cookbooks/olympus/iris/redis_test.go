package iris_olympus_cookbook

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func (t *IrisCookbookTestSuite) TestDeployRedis() {
	t.TestUploadRedis()
	//cdep := redisClusterDefinition.GenerateDeploymentRequest()

	//_, err := t.ZeusTestClient.DeployCluster(ctx, cdep)
	//t.Require().Nil(err)
}

func (t *IrisCookbookTestSuite) TestDestroyRedis() {
	d := zeus_req_types.TopologyDeployRequest{
		CloudCtxNs: redisCloudCtxNs,
	}
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, d)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *IrisCookbookTestSuite) TestUploadRedis() {
	_, rerr := redisClusterDefinition.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(rerr)
}

func (t *IrisCookbookTestSuite) TestCreateClusterClassRedis() {
	cd := redisClusterDefinition
	gcd := cd.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	err := gcd.CreateClusterClassDefinitions(context.Background(), t.ZeusTestClient)
	t.Require().Nil(err)
}
