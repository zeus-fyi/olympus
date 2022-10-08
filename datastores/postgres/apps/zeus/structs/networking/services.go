package networking

import (
	autogen_structs2 "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
)

type Service struct {
	KindDefinition        autogen_structs2.ChartComponentKinds
	ParentClassDefinition autogen_structs2.ChartSubcomponentParentClassTypes

	Metadata common.Metadata
	ServiceSpec
}

func NewService() Service {
	s := Service{}
	s.KindDefinition = autogen_structs2.ChartComponentKinds{
		ChartComponentKindName:   "Service",
		ChartComponentApiVersion: "v1",
	}
	s.ServiceSpec = NewServiceSpec()
	return s
}
