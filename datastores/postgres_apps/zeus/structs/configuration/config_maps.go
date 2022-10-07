package configuration

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
)

type ConfigMap struct {
	KindDefinition        autogen_structs.ChartComponentKinds
	ParentClassDefinition autogen_structs.ChartSubcomponentParentClassTypes
}

func NewConfigMap() ConfigMap {
	cm := ConfigMap{}
	cm.KindDefinition = autogen_structs.ChartComponentKinds{
		ChartComponentKindName:   "ConfigMap",
		ChartComponentApiVersion: "v1",
	}
	return cm
}
