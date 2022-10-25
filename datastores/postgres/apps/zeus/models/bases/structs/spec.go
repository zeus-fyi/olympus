package structs

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type SpecWorkload struct {
	autogen_bases.ChartSubcomponentParentClassTypes

	Replicas ChildClassSingleValue
	Selector Selector
}

func (s *SpecWorkload) GetReplicaCount32IntPtr() *int32 {
	return string_utils.ConvertStringTo32BitPtrInt(s.Replicas.ChartSubcomponentValue)
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
