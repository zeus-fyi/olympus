package configuration

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

type ConfigMap struct {
	KindDefinition        autogen_bases.ChartComponentResources
	ParentClassDefinition autogen_bases.ChartSubcomponentParentClassTypes
}

func NewConfigMap() ConfigMap {
	cm := ConfigMap{}
	cm.KindDefinition = autogen_bases.ChartComponentResources{
		ChartComponentKindName:   "ConfigMap",
		ChartComponentApiVersion: "v1",
	}
	return cm
}
