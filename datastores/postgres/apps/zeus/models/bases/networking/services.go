package networking

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
)

type Service struct {
	KindDefinition        autogen_bases.ChartComponentResources
	ParentClassDefinition autogen_bases.ChartSubcomponentParentClassTypes

	Metadata structs.Metadata
	ServiceSpec
}

func NewService() Service {
	s := Service{}
	s.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "Service",
		ChartComponentApiVersion: "v1",
	}
	s.ServiceSpec = NewServiceSpec()
	return s
}
