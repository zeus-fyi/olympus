package structs

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

func NewMetadataName() ChildClassSingleValue {
	return ChildClassSingleValue{
		ChartSubcomponentChildClassTypes: autogen_bases.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: "name",
		},
		ChartSubcomponentsChildValues: autogen_bases.ChartSubcomponentsChildValues{
			ChartSubcomponentKeyName:                       "name",
			ChartSubcomponentValue:                         "",
			ChartSubcomponentChildClassTypeID:              0,
			ChartSubcomponentChartPackageTemplateInjection: false,
		},
	}
}
