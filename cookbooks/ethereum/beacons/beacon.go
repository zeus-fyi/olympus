package beacon_cookbooks

import (
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

const (
	GenericBaseComponentsID = 1670202733617147904

	LocalEthereumBeaconClusterDefinitionID = 1670201797184939008

	ConsensusClientsBaseComponentsID = 1670202869405165056
	ExecClientsBaseComponentsID      = 1670202869402443776

	LighthouseSkeletonBaseID = 1670203661219772928
	GethSkeletonBaseID       = 1670203700436209920
)

var BeaconCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "ethereum", // set with your own namespace
	Env:           "production",
}
