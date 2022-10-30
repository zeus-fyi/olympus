package coreK8s

import clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/cluster"

type TopologyActionRequest struct {
	K8sRequest
	Action string

	clusters.Cluster
}

func (t *TopologyActionRequest) DeployTopology() {

	//chart := t.GetInfraChartPackage()

}
