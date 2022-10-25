package structs

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type ChildClassSingleValue struct {
	autogen_bases.ChartSubcomponentChildClassTypes
	autogen_bases.ChartSubcomponentsChildValues
}

func NewChildClassSingleValue(typeName string) ChildClassSingleValue {
	csv := ChildClassSingleValue{
		ChartSubcomponentChildClassTypes: autogen_bases.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: typeName,
		},
		ChartSubcomponentsChildValues: autogen_bases.ChartSubcomponentsChildValues{
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
		ChartSubcomponentChildClassTypes: autogen_bases.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: typeName,
		},
		ChartSubcomponentsChildValues: autogen_bases.ChartSubcomponentsChildValues{
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

func (c *ChildClassSingleValue) SetChildClassTypeIDs(id int) {
	c.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID = id
	c.ChartSubcomponentsChildValues.ChartSubcomponentChildClassTypeID = id
}

func (c *ChildClassSingleValue) GetChildClassTypeID() int {
	return c.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID
}

func (c *ChildClassSingleValue) GetParentClassTypeID() int {
	return c.ChartSubcomponentChildClassTypes.ChartSubcomponentParentClassTypeID
}

func (c *ChildClassSingleValue) GetChildValueTypeID() int {
	return c.ChartSubcomponentsChildValues.ChartSubcomponentChildClassTypeID
}
