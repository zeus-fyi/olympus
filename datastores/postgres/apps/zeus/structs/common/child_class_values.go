package common

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type ChildClassMultiValue struct {
	autogen_bases.ChartSubcomponentChildClassTypes
	Values []autogen_bases.ChartSubcomponentsChildValues
}

func NewChildClassAndValues(typeName string) ChildClassMultiValue {
	cvals := ChildClassMultiValue{
		ChartSubcomponentChildClassTypes: autogen_bases.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: typeName,
		},
		Values: []autogen_bases.ChartSubcomponentsChildValues{},
	}
	return cvals
}

func (c *ChildClassMultiValue) AddValues(values ...string) {
	if len(c.Values) <= 0 {
		c.Values = []autogen_bases.ChartSubcomponentsChildValues{}
	}
	for _, v := range values {
		c.Values = append(c.Values, autogen_bases.ChartSubcomponentsChildValues{ChartSubcomponentValue: v})
	}
}
