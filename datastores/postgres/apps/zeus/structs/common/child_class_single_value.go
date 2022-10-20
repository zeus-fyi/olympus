package common

import autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"

type ChildClassSingleValue struct {
	autogen_structs.ChartSubcomponentChildClassTypes
	autogen_structs.ChartSubcomponentsChildValues
}

func NewChildClassSingleValue(typeName string) ChildClassSingleValue {
	csv := ChildClassSingleValue{
		ChartSubcomponentChildClassTypes: autogen_structs.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: typeName,
		},
		ChartSubcomponentsChildValues: autogen_structs.ChartSubcomponentsChildValues{
			ChartSubcomponentChildClassTypeID:              0,
			ChartSubcomponentChartPackageTemplateInjection: false,
			ChartSubcomponentKeyName:                       typeName,
			ChartSubcomponentValue:                         "",
		},
	}
	return csv
}

func NewInitChildClassSingleValue(typeName, value string) ChildClassSingleValue {
	csv := ChildClassSingleValue{
		ChartSubcomponentChildClassTypes: autogen_structs.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: typeName,
		},
		ChartSubcomponentsChildValues: autogen_structs.ChartSubcomponentsChildValues{
			ChartSubcomponentChildClassTypeID:              0,
			ChartSubcomponentChartPackageTemplateInjection: false,
			ChartSubcomponentKeyName:                       typeName,
			ChartSubcomponentValue:                         value,
		},
	}
	return csv
}

func (c *ChildClassSingleValue) SetValue(v string) {
	c.ChartSubcomponentValue = v
}
