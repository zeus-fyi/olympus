package common

import autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"

type MetadataName struct {
	autogen_structs.ChartSubcomponentChildClassTypes
	autogen_structs.ChartSubcomponentsChildValues
}

func NewMetadataName() MetadataName {
	return MetadataName{
		autogen_structs.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: "name",
		},
		autogen_structs.ChartSubcomponentsChildValues{
			ChartSubcomponentChildClassTypeID:              0,
			ChartSubcomponentChartPackageTemplateInjection: false,
			ChartSubcomponentKeyName:                       "name",
			ChartSubcomponentValue:                         "",
		},
	}
}

func (mn *MetadataName) AddNameValue(name string) {
	mn.ChartSubcomponentValue = name
}
