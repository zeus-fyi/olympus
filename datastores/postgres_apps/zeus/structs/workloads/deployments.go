package workloads

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
)

type Deployment struct {
	ClassDefinition autogen_structs.ChartComponentKinds
	Metadata        common.Metadata
}

func NewDeployment() Deployment {
	d := Deployment{}
	d.ClassDefinition = autogen_structs.ChartComponentKinds{
		ChartComponentKindName:   "Deployment",
		ChartComponentApiVersion: "apps/v1",
	}
	return d
}
