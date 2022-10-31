package clusters

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases"
	topology2 "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology"
)

type Cluster struct {
	topology2.Topology
	Base       bases.Base
	Components topology2.Dependencies
}

func NewCluster() Cluster {
	cl := Cluster{topology2.NewClusterTopology(), bases.NewBase(), []topology2.Component{}}
	return cl
}

func (c *Cluster) GetInfraChartPackage() charts.Chart {
	return c.Base.Chart
}
