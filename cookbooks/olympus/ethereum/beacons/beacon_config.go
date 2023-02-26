package olympus_beacon_cookbooks

import (
	zeus_topology_config_drivers "github.com/zeus-fyi/zeus/pkg/zeus/workload_config_drivers"
	v1 "k8s.io/api/core/v1"
)

// TODO
var (
	consensusClientStsEphemeralCfg = zeus_topology_config_drivers.StatefulSetDriver{
		ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
			"zeus-consensus-client": {
				Container: v1.Container{
					Name:  "zeus-consensus-client",
					Image: "sigp/lighthouse:capella",
				},
			},
		}}
	consensusClientStsEphemeralCfgDriver = zeus_topology_config_drivers.TopologyConfigDriver{
		StatefulSetDriver: &consensusClientStsEphemeralCfg,
	}
	execClientStsEphemeralCfg = zeus_topology_config_drivers.StatefulSetDriver{
		ContainerDrivers: map[string]zeus_topology_config_drivers.ContainerDriver{
			"zeus-exec-client": {
				Container: v1.Container{
					Name:  "zeus-exec-client",
					Image: "ethpandaops/geth:master",
				},
			},
		}}
	execClientStsEphemeralCfgDriver = zeus_topology_config_drivers.TopologyConfigDriver{
		StatefulSetDriver: &execClientStsEphemeralCfg,
	}
)
