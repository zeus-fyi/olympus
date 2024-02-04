package zeus_v1_clusters_api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

var ctx = context.Background()

type KubeConfigRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *KubeConfigRequestTestSuite) TestKubeConfigUpload() {
	t.Eg.POST("/kubeconfig", CreateOrUpdateKubeConfigsHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

}

func TestKubeConfigRequestTestSuite(t *testing.T) {
	suite.Run(t, new(KubeConfigRequestTestSuite))
}
