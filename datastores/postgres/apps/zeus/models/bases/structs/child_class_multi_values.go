package structs

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type ChildClassMultiValue struct {
	autogen_bases.ChartSubcomponentChildClassTypes
	Values []autogen_bases.ChartSubcomponentsChildValues
}

func NewChildClassMultiValues(typeName string) ChildClassMultiValue {
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

func (c *ChildClassMultiValue) AddValuesAndUniqueChildID(id int, values map[string]string) {
	if len(c.Values) <= 0 {
		c.Values = []autogen_bases.ChartSubcomponentsChildValues{}
	}
	for k, v := range values {

		cv := autogen_bases.ChartSubcomponentsChildValues{ChartSubcomponentKeyName: k, ChartSubcomponentValue: v}
		cv.ChartSubcomponentChildClassTypeID = id
		c.Values = append(c.Values, cv)
	}
}

func (c *ChildClassMultiValue) AddValues(values map[string]string) {
	if len(c.Values) <= 0 {
		c.Values = []autogen_bases.ChartSubcomponentsChildValues{}
	}
	for k, v := range values {
		cv := autogen_bases.ChartSubcomponentsChildValues{ChartSubcomponentKeyName: k, ChartSubcomponentValue: v}
		c.Values = append(c.Values, cv)
	}
}

func (c *ChildClassMultiValue) SetParentClassTypeID(id int) {
	c.ChartSubcomponentChildClassTypes.ChartSubcomponentParentClassTypeID = id
}

func (c *ChildClassMultiValue) GetMultiValueChildClassTypeID() int {
	return c.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID
}

func (c *ChildClassMultiValue) SetMultiValueChildClassTypeIDs(id int) {
	c.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID = id

	for i, _ := range c.Values {
		c.Values[i].ChartSubcomponentChildClassTypeID = id
	}
}
