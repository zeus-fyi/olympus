package common

import autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"

type ChildClassAndValues struct {
	autogen_structs.ChartSubcomponentChildClassTypes
	Values []autogen_structs.ChartSubcomponentsChildValues
}

func NewChildClassAndValues(typeName string) ChildClassAndValues {
	cvals := ChildClassAndValues{
		ChartSubcomponentChildClassTypes: autogen_structs.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: typeName,
		},
		Values: []autogen_structs.ChartSubcomponentsChildValues{},
	}
	return cvals
}

func (c *ChildClassAndValues) AddValues(values ...autogen_structs.ChartSubcomponentsChildValues) {
	if len(c.Values) <= 0 {
		c.Values = []autogen_structs.ChartSubcomponentsChildValues{}
	}
	for _, v := range values {
		c.Values = append(c.Values, v)
	}
}
