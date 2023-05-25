package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type DeploymentsTestSuite struct {
	K8TestSuite
}

func (d *DeploymentsTestSuite) TestCreateDeployment() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "p2p-crawler"}
	dep, err := d.K.GetDeployment(ctx, kns, "zeus-p2p-crawler", nil)
	d.Require().Nil(err)
	d.Require().NotEmpty(dep)

	dep.ResourceVersion = ""
	_, err = d.K.CreateDeployment(ctx, kns, dep, nil)
	d.Require().Nil(err)
}

func (d *DeploymentsTestSuite) TestGetDeployment() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "", Region: "", Context: "data", Namespace: "eth-indexer"}
	dep, err := d.K.GetDeployment(ctx, kns, "eth-indexer", nil)
	d.Require().Nil(err)
	d.Require().NotEmpty(dep)
}

func TestDeploymentsTestSuite(t *testing.T) {
	suite.Run(t, new(DeploymentsTestSuite))
}
