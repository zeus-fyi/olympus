package networking

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
)

type Service struct {
	ClassDefinition autogen_structs.ChartComponentKinds
}

func NewService() Service {
	s := Service{}
	s.ClassDefinition = autogen_structs.ChartComponentKinds{
		ChartComponentKindName:   "Service",
		ChartComponentApiVersion: "v1",
	}
	return s
}
