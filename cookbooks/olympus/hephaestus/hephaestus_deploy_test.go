package hephaestus_olympus_cookbook

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
)

func (t *HephaestusCookbookTestSuite) TestDeploy() {
	_, rerr := HephaestusClusterDefinition.UploadChartsFromClusterDefinition(ctx, t.ZeusTestClient, true)
	t.Require().Nil(rerr)
	cdep := HephaestusClusterDefinition.GenerateDeploymentRequest()

	_, err := t.ZeusTestClient.DeployCluster(ctx, cdep)
	t.Require().Nil(err)
}

func (t *HephaestusCookbookTestSuite) TestClusterDestroy() {
	d := zeus_req_types.TopologyDeployRequest{
		CloudCtxNs: HephaestusCloudCtxNs,
	}
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, d)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *HephaestusCookbookTestSuite) TestCreateClusterClass() {
	cd := HephaestusClusterDefinition
	gcd := cd.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	err := gcd.CreateClusterClassDefinitions(context.Background(), t.ZeusTestClient)
	t.Require().Nil(err)
}
