package olympus_beacon

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

// set your own topologyID here after uploading a chart workload
var deployExecKnsReq = zeus_req_types.TopologyDeployRequest{
	TopologyID: 1669168567793564000,
	CloudCtxNs: topCloudCtxNs,
}

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var execChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./olympus/beacon/exec_client",
	DirOut:      "./olympus/outputs",
	FnIn:        "execClientAthena", // filename for your gzip workload
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

	chart := newUploadChart(basePath.FnIn)
	_, err := t.ZeusTestClient.UploadChart(ctx, basePath, chart)
	t.Require().Nil(err)

	tar := zeus_req_types.TopologyRequest{TopologyID: deployExecKnsReq.TopologyID}
	resp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	// prints the chart output for inspection
	err = resp.PrintWorkload(basePath)
}
