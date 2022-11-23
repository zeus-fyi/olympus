package ethereum_beacon_cookbook

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var consensusChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./ethereum/beacon/consensus_client",
	DirOut:      "./ethereum/outputs",
	FnIn:        "consensusClientHercules", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}

// DirOut is where it will write a copy of the chart you uploaded, which helps verify the workload is correct
var execChartPath = filepaths.Path{
	PackageName: "",
	DirIn:       "./ethereum/beacon/exec_client",
	DirOut:      "./ethereum/outputs",
	FnIn:        "execClientHercules", // filename for your gzip workload
	FnOut:       "",
	Env:         "",
	FilterFiles: string_utils.FilterOpts{},
}

func (t *ZeusEthereumBeaconTestSuite) TestBeaconUpload() {
	basePath := consensusChartPath

	// derived
	chart := newUploadChart(basePath.FnIn)
	uploadResp, err := t.ZeusTestClient.UploadChart(ctx, basePath, chart)
	t.Require().Nil(err)

	tar := zeus_req_types.TopologyRequest{TopologyID: uploadResp.TopologyID}
	resp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	// prints the chart output for inspection
	err = resp.PrintWorkload(basePath)
	basePath = execChartPath

	chart = newUploadChart(basePath.FnIn)
	uploadResp, err = t.ZeusTestClient.UploadChart(ctx, basePath, chart)
	t.Require().Nil(err)

	tar = zeus_req_types.TopologyRequest{TopologyID: uploadResp.TopologyID}
	resp, err = t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

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
