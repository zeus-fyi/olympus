package create_infra

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases/infra"
)

type Infrastructure struct {
	infra.Infrastructure
}

func NewCreateInfrastructure() Infrastructure {
	inf := infra.NewInfrastructure()
	i := Infrastructure{inf}
	return i
}
