package zeus_apps

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	zeus_pods_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/pods"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

func (t *ZeusAppsTestSuite) TestGenericExecChartUploadAndRead() {
	basePath := execChartPath

	// derived
	chart := newUploadChart(basePath.FnIn)
	resp, err := t.ZeusTestClient.UploadChart(ctx, basePath, chart)
	t.Require().Nil(err)
	t.TestReadDemoChart(resp.TopologyID, basePath)
}

func (t *ZeusAppsTestSuite) TestGenericExecChartDeploy() {
	_, err := t.ZeusTestClient.Deploy(ctx, deployKnsReq)
	t.Require().Nil(err)
}

// set your own topologyID here after uploading a chart workload
var deployConsensusKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1668825071425428000,
	CloudCtxNs: topCloudCtxNs,
}

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

	deployKnsReq.Namespace = "ethereum"
	par := zeus_pods_reqs.PodActionRequest{
		TopologyDeployRequest: deployConsensusKnsReq,
		Action:                zeus_pods_reqs.DeleteAllPods,
		PodName:               "zeus-lighthouse-0",
	}
	err = t.ZeusTestClient.DeletePods(ctx, par)
	t.Require().Nil(err)
}

// set your own topologyID here after uploading a chart workload
var deployKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1668887498372169000,
	CloudCtxNs: topCloudCtxNs,
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

var topCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "ethereum", // set with your own namespace
	Env:           "dev",
}

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var consensusChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./mocks/kubernetes_apps/beacon/consensus_client",
	DirOut:      "./outputs/consensus_client_out",
	FnIn:        "consensus_client", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var execChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./mocks/kubernetes_apps/beacon/exec_client_generic",
	DirOut:      "./outputs/generic_exec_chart_out",
	FnIn:        "exec", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}
