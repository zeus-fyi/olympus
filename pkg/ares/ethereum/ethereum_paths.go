package ethereum

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func ExecClientPath() filepaths.Path {
	var execClientPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./beacon/exec_client",
		DirOut:      "./beacon_out/exec_client_out/gzip",
		FnIn:        "exec_client",
		FnOut:       "",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	return execClientPath
}

func ExecClientReadChartThenWritePath() filepaths.Path {
	var execClientPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./beacon/exec_client",
		DirOut:      "./beacon_out/exec_client_out/read_chart",
		FnIn:        "",
		FnOut:       "",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	return execClientPath
}

func ConsensusClientPath() filepaths.Path {
	var consensusClientPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./beacon/consensus_client",
		DirOut:      "./beacon_out/consensus_client_out/gzip",
		FnIn:        "consensus_client",
		FnOut:       "",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	return consensusClientPath
}

func ConsensusReadChartThenWritePath() filepaths.Path {
	var consensusClientPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./beacon/consensus_client",
		DirOut:      "./beacon_out/consensus_client_out/read_chart",
		FnIn:        "",
		FnOut:       "",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	return consensusClientPath
}
