package ethereum

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

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

func ConsensusReadChartThenWritePath() structs.Path {
	var consensusClientPath = structs.Path{
		PackageName: "",
		DirIn:       "./beacon/consensus_client",
		DirOut:      "./beacon_out/consensus_client_out/read_chart",
		Fn:          "",
		FnOut:       "",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	return consensusClientPath
}
