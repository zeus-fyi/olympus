package configuration

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
)

type ConfigMap struct {
	ClassDefinition autogen_structs.ChartComponentKinds
}

func NewConfigMap() ConfigMap {
	cm := ConfigMap{}
	cm.ClassDefinition = autogen_structs.ChartComponentKinds{
		ChartComponentKindName:   "ConfigMap",
		ChartComponentApiVersion: "v1",
	}
	return cm
}
