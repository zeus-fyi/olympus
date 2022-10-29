package infra

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
)

type Infrastructure struct {
	autogen_bases.TopologyInfrastructureComponents
	charts.Chart
}

func NewInfrastructure() Infrastructure {
	inf := Infrastructure{autogen_bases.TopologyInfrastructureComponents{
		TopologyID:     0,
		ChartPackageID: 0,
	}, charts.NewChart()}
	return inf
}

func (i *Infrastructure) AddChart(chart charts.Chart) {
	i.ChartPackageID = chart.ChartPackageID
	i.Chart = chart
}
