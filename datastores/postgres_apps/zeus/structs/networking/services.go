package networking

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
)

type Service struct {
	KindDefinition        autogen_structs.ChartComponentKinds
	ParentClassDefinition autogen_structs.ChartSubcomponentParentClassTypes

	Metadata common.Metadata
	ServiceSpec
}

func NewService() Service {
	s := Service{}
	s.KindDefinition = autogen_structs.ChartComponentKinds{
		ChartComponentKindName:   "Service",
		ChartComponentApiVersion: "v1",
	}
	s.ServiceSpec = NewServiceSpec()
	return s
}
