package olympus_beacon

import zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"

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
	_, err = t.ZeusTestClient.DeletePods(ctx, par)
	t.Require().Nil(err)
}
