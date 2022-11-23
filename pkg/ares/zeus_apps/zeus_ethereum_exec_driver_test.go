package zeus_apps

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

// set your own topologyID here after uploading a chart workload
var deployExecKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1668887498372169000,
	CloudCtxNs: topCloudCtxNs,
}

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var execChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./mocks/kubernetes_apps/beacon/exec_client",
	DirOut:      "./mocks/kubernetes_apps/beacon/exec_client_out",
	FnIn:        "exec", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}

func (t *ZeusAppsTestSuite) TestExecChartDeploy() {
	_, err := t.ZeusTestClient.Deploy(ctx, deployExecKnsReq)
	t.Require().Nil(err)
}

func (t *ZeusAppsTestSuite) TestExecClientUpload() {
	basePath := execChartPath

	// derived
	chart := newUploadChart(basePath.FnIn)
	resp, err := t.ZeusTestClient.UploadChart(ctx, basePath, chart)
	t.Require().Nil(err)
	t.TestReadDemoChart(resp.TopologyID, basePath)
}
