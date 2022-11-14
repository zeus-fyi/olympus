package ethereum

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	create_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"
)

var ts chronos.Chronos

func ConsensusClientChartUploadRequest() create_infra.TopologyCreateRequest {
	uploadChart := create_infra.TopologyCreateRequest{
		TopologyName:     "consensusClient",
		ChartName:        "lighthouse",
		ChartDescription: "lighthouse client",
		Version:          fmt.Sprintf("v0.0.%d", ts.UnixTimeStampNow()),
	}
	return uploadChart
}

func ExecClientChartUploadRequest() create_infra.TopologyCreateRequest {
	uploadChart := create_infra.TopologyCreateRequest{
		TopologyName:     "execClient",
		ChartName:        "geth",
		ChartDescription: "geth client",
		Version:          fmt.Sprintf("v0.0.%d", ts.UnixTimeStampNow()),
	}
	return uploadChart
}

func ChangeDirToAresEthereumDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}
