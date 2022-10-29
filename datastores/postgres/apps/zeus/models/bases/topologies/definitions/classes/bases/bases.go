package bases

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases/config"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases/infra"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/bases/skeleton"
)

type Base struct {
	infra.Infrastructure
	*config.Configuration
	*skeleton.Skeleton
}

func NewBase() Base {
	b := Base{}
	return b
}
