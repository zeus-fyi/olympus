package common

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
)

type Metadata struct {
	MetadataClassDefinition autogen_structs.ChartSubcomponentChildClassTypes
	Name                    string
	Annotations             ChildValuesSlice
	Labels                  ChildValuesSlice
}

func NewMetadata() Metadata {
	m := Metadata{}
	m.MetadataClassDefinition = autogen_structs.ChartSubcomponentChildClassTypes{
		ChartSubcomponentParentClassTypeID:  0,
		ChartSubcomponentChildClassTypeID:   0,
		ChartSubcomponentChildClassTypeName: "metadata",
	}

	// (volumes, nodeSelector, affinity, tolerations, etc)
	m.Annotations = ChildValuesSlice{}
	m.Labels = ChildValuesSlice{}
	return m
}
