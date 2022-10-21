package configuration

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

type ConfigMap struct {
	KindDefinition        autogen_structs.ChartComponentResources
	ParentClassDefinition autogen_structs.ChartSubcomponentParentClassTypes
}

func NewConfigMap() ConfigMap {
	cm := ConfigMap{}
	cm.KindDefinition = autogen_structs.ChartComponentResources{
		ChartComponentKindName:   "ConfigMap",
		ChartComponentApiVersion: "v1",
	}
	return cm
}
