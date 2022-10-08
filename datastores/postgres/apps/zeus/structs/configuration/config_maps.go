package configuration

import (
	autogen_structs2 "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

type ConfigMap struct {
	KindDefinition        autogen_structs2.ChartComponentKinds
	ParentClassDefinition autogen_structs2.ChartSubcomponentParentClassTypes
}

func NewConfigMap() ConfigMap {
	cm := ConfigMap{}
	cm.KindDefinition = autogen_structs2.ChartComponentKinds{
		ChartComponentKindName:   "ConfigMap",
		ChartComponentApiVersion: "v1",
	}
	return cm
}
