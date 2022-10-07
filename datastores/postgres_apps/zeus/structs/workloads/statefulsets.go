package workloads

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/networking"
)

type StatefulSet struct {
	ClassDefinition autogen_structs.ChartComponentKinds
	Metadata        common.Metadata

	ServiceDefinition networking.Service
}

func NewStatefulSet() StatefulSet {
	s := StatefulSet{}
	s.ClassDefinition = autogen_structs.ChartComponentKinds{
		ChartComponentKindName:   "StatefulSet",
		ChartComponentApiVersion: "apps/v1",
	}

	return s
}
