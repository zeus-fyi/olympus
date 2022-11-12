package ethereum

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
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

func ConsensusClientPath() structs.Path {
	var consensusClientPath = structs.Path{
		PackageName: "",
		DirIn:       "./beacon/consensus_client",
		DirOut:      "./beacon_out/consensus_client_out/gzip",
		Fn:          "consensus_client",
		FnOut:       "",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	return consensusClientPath
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
