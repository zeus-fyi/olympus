package internal_reqs

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"

type InternalSecretsCopyFromTo struct {
	SecretNames []string              `json:"secretNames"`
	FromKns     kns.TopologyKubeCtxNs `json:"fromKns"`
	ToKns       kns.TopologyKubeCtxNs `json:"toKns"`
}
