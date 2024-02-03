package zeus_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type DeploymentsTestSuite struct {
	K8TestSuite
}

func (d *DeploymentsTestSuite) TestCreateDeployment() {
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "p2p-crawler"}
	dep, err := d.K.GetDeployment(ctx, kns, "zeus-p2p-crawler", nil)
	d.Require().Nil(err)
	d.Require().NotEmpty(dep)

	dep.ResourceVersion = ""
	_, err = d.K.CreateDeployment(ctx, kns, dep, nil)
	d.Require().Nil(err)
}

func (d *DeploymentsTestSuite) TestGetDeployment() {
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "", Region: "", Context: "data", Namespace: "eth-indexer"}
	dep, err := d.K.GetDeployment(ctx, kns, "eth-indexer", nil)
	d.Require().Nil(err)
	d.Require().NotEmpty(dep)
}

func (d *DeploymentsTestSuite) TestRolloutRestartDeploymentZeus() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "kubernetes-admin@zeusfyi", Namespace: "zeus"}
	dep, err := d.K.RolloutRestartDeployment(ctx, kns, "zeus", nil)
	d.Require().Nil(err)
	d.Require().NotEmpty(dep)

	kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "zeus"}
	dep, err = d.K.RolloutRestartDeployment(ctx, kns, "zeus", nil)
	d.Require().Nil(err)
	d.Require().NotEmpty(dep)
}

func (d *DeploymentsTestSuite) TestRolloutRestartDeploymentWebApp() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "kubernetes-admin@zeusfyi", Namespace: "hestia"}
	dep, err := d.K.RolloutRestartDeployment(ctx, kns, "hestia", nil)
	d.Require().Nil(err)
	d.Require().NotEmpty(dep)

	dep, err = d.K.RolloutRestartDeployment(ctx, kns, "zeus-cloud", nil)
	d.Require().Nil(err)
	d.Require().NotEmpty(dep)

	kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "hestia"}
	dep, err = d.K.RolloutRestartDeployment(ctx, kns, "hestia", nil)
	d.Require().Nil(err)
}

func (d *DeploymentsTestSuite) TestRolloutRestartDeploymentHestia() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "kubernetes-admin@zeusfyi", Namespace: "hestia"}
	dep, err := d.K.RolloutRestartDeployment(ctx, kns, "hestia", nil)
	d.Require().Nil(err)
	d.Require().NotEmpty(dep)

	kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "hestia"}
	dep, err = d.K.RolloutRestartDeployment(ctx, kns, "hestia", nil)
	d.Require().Nil(err)
}
func (d *DeploymentsTestSuite) TestRolloutRestartDeploymentZeusCloud() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "kubernetes-admin@zeusfyi", Namespace: "hestia"}
	dep, err := d.K.RolloutRestartDeployment(ctx, kns, "zeus-cloud", nil)
	d.Require().Nil(err)
	d.Require().NotEmpty(dep)
}

func (d *DeploymentsTestSuite) TestRolloutRestartDeploymentIris() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "kubernetes-admin@zeusfyi", Namespace: "iris"}
	dep, err := d.K.RolloutRestartDeployment(ctx, kns, "iris", nil)
	d.Require().Nil(err)
	d.Require().NotEmpty(dep)
}

func TestDeploymentsTestSuite(t *testing.T) {
	suite.Run(t, new(DeploymentsTestSuite))
}
