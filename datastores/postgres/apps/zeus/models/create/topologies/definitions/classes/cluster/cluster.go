package create_clusters

import clusters "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/cluster"

type Cluster struct {
	clusters.Cluster
}

func NewCreateCluster() Cluster {
	c := Cluster{
		clusters.NewCluster(),
	}
	return c
}
