package common

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

type Selector struct {
	ClassDefinition autogen_structs.ChartSubcomponentChildClassTypes
	MatchLabels     ChildValuesSlice
}

func NewSelector() Selector {
	s := Selector{}
	s.ClassDefinition = autogen_structs.ChartSubcomponentChildClassTypes{
		ChartSubcomponentParentClassTypeID:  0,
		ChartSubcomponentChildClassTypeID:   0,
		ChartSubcomponentChildClassTypeName: "selector",
	}

	// (volumes, nodeSelector, affinity, tolerations, etc)
	s.MatchLabels = ChildValuesSlice{}
	return s
}
