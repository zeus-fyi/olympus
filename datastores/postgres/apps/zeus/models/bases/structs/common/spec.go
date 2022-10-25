package common

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

type SpecWorkload struct {
	autogen_bases.ChartSubcomponentParentClassTypes

	Replicas ChildClassSingleValue
	Selector Selector
}

func NewSpecWorkload() SpecWorkload {
	pc := autogen_bases.ChartSubcomponentParentClassTypes{
		ChartPackageID:                       0,
		ChartComponentResourceID:             0,
		ChartSubcomponentParentClassTypeID:   0,
		ChartSubcomponentParentClassTypeName: "Spec",
	}
	sw := SpecWorkload{pc, NewInitChildClassSingleValue("replicas", "0"), NewSelector()}
	return sw
}

func (s *SpecWorkload) SetParentClassTypeIDs(id int) {
	s.ChartSubcomponentParentClassTypeID = id
	s.Selector.MatchLabels.ChartSubcomponentParentClassTypeID = id
	s.Replicas.ChartSubcomponentParentClassTypeID = id
}
