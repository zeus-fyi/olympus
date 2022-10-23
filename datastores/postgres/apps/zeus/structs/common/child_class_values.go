package common

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type ChildClassAndValues struct {
	autogen_bases.ChartSubcomponentChildClassTypes
	Values []autogen_bases.ChartSubcomponentsChildValues
}

func NewChildClassAndValues(typeName string) ChildClassAndValues {
	cvals := ChildClassAndValues{
		ChartSubcomponentChildClassTypes: autogen_bases.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: typeName,
		},
		Values: []autogen_bases.ChartSubcomponentsChildValues{},
	}
	return cvals
}

func (c *ChildClassAndValues) AddValues(values ...autogen_bases.ChartSubcomponentsChildValues) {
	if len(c.Values) <= 0 {
		c.Values = []autogen_bases.ChartSubcomponentsChildValues{}
	}
	for _, v := range values {
		c.Values = append(c.Values, v)
	}
}
