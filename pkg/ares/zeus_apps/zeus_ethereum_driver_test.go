package zeus_apps

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
)

func (t *ZeusAppsTestSuite) TestRead() {
	_, err := t.ZeusTestClient.ReadTopologies(ctx)
	t.Require().Nil(err)
}

func (t *ZeusAppsTestSuite) TestReplaceCm() {
	basePath := consensusChartPath
	basePath.DirIn += "/alt_config"
	resp, err := t.ZeusTestClient.DeployReplace(ctx, basePath, deployConsensusKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	deployConsensusKnsReq.Namespace = "beacon"
	par := zeus_pods_reqs.PodActionRequest{
		TopologyDeployRequest: deployConsensusKnsReq,
		Action:                zeus_pods_reqs.DeleteAllPods,
		PodName:               "zeus-lighthouse-0",
	}
	err = t.ZeusTestClient.DeletePods(ctx, par)
	t.Require().Nil(err)
}

// chart workload metadata
func newUploadChart(name string) zeus_req_types.TopologyCreateRequest {
	return zeus_req_types.TopologyCreateRequest{
		TopologyName:     name,
		ChartName:        name,
		ChartDescription: name,
		Version:          fmt.Sprintf("v0.0.%d", time.Now().Unix()),
	}
}
