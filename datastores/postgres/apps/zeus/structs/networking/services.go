package networking

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
)

type Service struct {
	KindDefinition        autogen_structs.ChartComponentResources
	ParentClassDefinition autogen_structs.ChartSubcomponentParentClassTypes

	Metadata common.Metadata
	ServiceSpec
}

func NewService() Service {
	s := Service{}
	s.KindDefinition = autogen_structs.ChartComponentResources{
		ChartComponentKindName:   "Service",
		ChartComponentApiVersion: "v1",
	}
	s.ServiceSpec = NewServiceSpec()
	return s
}
