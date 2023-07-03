package athena_workloads

import (
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type WorkloadInfo struct {
	WorkloadType      string // eg, validatorClient
	ClientName        string // eg. lighthouse, geth
	ProtocolNetworkID int    // eg. mainnet
	ReplicaCountNum   int    // eg. stateful set ordinal index
	zeus_common_types.CloudCtxNs
	DataDir filepaths.Path
}
