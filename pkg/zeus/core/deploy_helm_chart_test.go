package zeus_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type HelmChartTestSuite struct {
	K8TestSuite
}

func (d *DeploymentsTestSuite) TestDeployHelmChartRelease() {
	ctxName := "arn:aws:eks:us-west-1:480391564655:cluster"
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "aws", Region: "us-west-1", Context: ctxName, Namespace: "observability"}

	err := d.K.DeployHelm(ctx, kns)
	d.Require().Nil(err)
}

// helm install [RELEASE_NAME] prometheus-community/kube-prometheus-stack

func TestHelmChartTestSuite(t *testing.T) {
	suite.Run(t, new(HelmChartTestSuite))
}
