package common_conversions

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

type ParentChildDB struct {
	PCGroupMap map[string][]PC
}

type ParentClassTypesDB struct {
	ChartSubcomponentParentClassTypeID   int    `json:"chart_subcomponent_parent_class_type_id"`
	ChartSubcomponentParentClassTypeName string `json:"chart_subcomponent_parent_class_type_name"`
}
type PC struct {
	ParentClassTypesDB                             `json:"chart_subcomponent_parent_class_types"` // chart_subcomponent_parent_class_types
	autogen_bases.ChartSubcomponentChildClassTypes `json:"chart_subcomponent_child_class_types"`  //chart_subcomponent_child_class_types
	autogen_bases.ChartSubcomponentsChildValues    `json:"chart_subcomponents_child_values"`      // chart_subcomponents_child_values
}

func NewPCGroupMap() ParentChildDB {
	m := make(map[string][]PC)
	return ParentChildDB{m}
}

func (pc *ParentChildDB) AppendPCElementToMap(keyName string, pcElement PC) {
	tmp := pc.PCGroupMap[keyName]
	tmp = append(tmp, pcElement)
	pc.PCGroupMap[keyName] = tmp
}
