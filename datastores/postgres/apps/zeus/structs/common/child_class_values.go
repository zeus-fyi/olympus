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

func (c *ChildClassMultiValue) AddValues(values map[string]string) {
	if len(c.Values) <= 0 {
		c.Values = []autogen_bases.ChartSubcomponentsChildValues{}
	}
	for k, v := range values {
		c.Values = append(c.Values, autogen_bases.ChartSubcomponentsChildValues{ChartSubcomponentKeyName: k, ChartSubcomponentValue: v})
	}
}

func (c *ChildClassMultiValue) GetClassTypeID() int {
	return c.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID
}

func (c *ChildClassMultiValue) SetClassTypeID(id int) {
	c.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID = id

	for _, v := range c.Values {
		v.ChartSubcomponentChildClassTypeID = id
	}
}
