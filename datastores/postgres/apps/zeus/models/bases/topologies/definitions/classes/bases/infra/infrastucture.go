package infra

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
)

type Infrastructure struct {
	charts.Chart
}

func NewInfrastructure() Infrastructure {
	inf := Infrastructure{charts.NewChart()}
	return inf
}
